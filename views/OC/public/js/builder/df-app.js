// Drawflow integration wrapper for goFlow
// Replaces previous custom app.js with Drawflow-based editor

let editor = null;
let templates = [];

function initDrawflow() {
  const el = document.getElementById('drawflow');
  editor = new Drawflow(el);
  // Basic config
  editor.start();
  // Optional: allow zoom with Ctrl+wheel (default handled in lib)
  // Use UUIDs so new nodes get stable non-numeric ids
  editor.useuuid = true;

  // Properties panel
  const props = document.getElementById('props');
  const renderProps = (id) => {
    try {
      const node = editor.getNodeFromId(id);
      const type = node?.data?.type || node?.name || '';
      const desc = node?.data?.description || '';
      const args = node?.data?.args || {};
      const outputs = node?.data?.outputs || [];
      const esc = (s) => String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;');
      const rows = Object.keys(args).map(k => `
        <div class="flex gap-1 items-center">
          <input class="border px-1 py-0.5 flex-1" data-arg-key value="${esc(k)}" />
          <input class="border px-1 py-0.5 flex-1" data-arg-val value="${esc(args[k])}" />
          <button class="text-xs border px-2 py-0.5" data-arg-del>Del</button>
        </div>`).join('');
      const outs = outputs.map((label, i) => `
        <div class="flex gap-1 items-center">
          <span class="text-xs text-gray-500 w-10">out${i+1}</span>
          <input class="border px-1 py-0.5 flex-1" data-out-idx="${i}" value="${esc(label)}" />
        </div>`).join('');
      props.innerHTML = `
        <div class="space-y-3">
          <div class="text-sm font-semibold">${esc(type)}</div>
          <div class="text-xs text-gray-500">${esc(desc)}</div>

          <div>
            <div class="text-xs font-semibold mb-1">Arguments</div>
            <div id="argRows" class="space-y-1">${rows || '<div class=\'text-xs text-gray-400\'>No args</div>'}</div>
            <div class="mt-1 flex gap-1">
              <button id="addArg" class="text-xs border px-2 py-0.5">Add arg</button>
              <button id="saveArgs" class="text-xs border px-2 py-0.5">Save args</button>
            </div>
          </div>

          <div>
            <div class="text-xs font-semibold mb-1">Outputs</div>
            <div id="outRows" class="space-y-1">${outs || '<div class=\'text-xs text-gray-400\'>No outputs</div>'}</div>
            <div class="mt-1 flex gap-1">
              <button id="addOut" class="text-xs border px-2 py-0.5">Add output</button>
              <button id="removeOut" class="text-xs border px-2 py-0.5">Remove last</button>
              <button id="saveOuts" class="text-xs border px-2 py-0.5">Save labels</button>
            </div>
          </div>
        </div>`;

      // Wire buttons
      props.querySelector('#addArg')?.addEventListener('click', () => {
        const c = props.querySelector('#argRows');
        const div = document.createElement('div');
        div.className = 'flex gap-1 items-center';
        div.innerHTML = `<input class="border px-1 py-0.5 flex-1" data-arg-key placeholder="key" />
                         <input class="border px-1 py-0.5 flex-1" data-arg-val placeholder="value" />
                         <button class="text-xs border px-2 py-0.5" data-arg-del>Del</button>`;
        c.appendChild(div);
      });
      props.addEventListener('click', (e) => {
        const t = e.target;
        if (t && t.matches('[data-arg-del]')) {
          t.parentElement.remove();
        }
      });
      props.querySelector('#saveArgs')?.addEventListener('click', () => {
        const rows = [...props.querySelectorAll('#argRows > div')];
        const next = {};
        rows.forEach(r => {
          const k = r.querySelector('[data-arg-key]')?.value?.trim();
          const v = r.querySelector('[data-arg-val]')?.value ?? '';
          if (k) next[k] = v;
        });
        node.data.args = next;
        // Optional: update node HTML with a small args preview
        updateNodeHtml(node.id);
      });
      props.querySelector('#addOut')?.addEventListener('click', () => {
        editor.addNodeOutput(node.id);
        const outs = node.data.outputs || [];
        outs.push('out'+(outs.length+1));
        node.data.outputs = outs;
        renderProps(node.id);
      });
      props.querySelector('#removeOut')?.addEventListener('click', () => {
        // remove last output if exists and no connections
        const outs = node.data.outputs || [];
        if (outs.length === 0) return;
        // Drawflow only supports removing last port via removeNodeOutput(id)
        editor.removeNodeOutput(node.id);
        outs.pop();
        node.data.outputs = outs;
        renderProps(node.id);
      });
      props.querySelector('#saveOuts')?.addEventListener('click', () => {
        const inputs = [...props.querySelectorAll('[data-out-idx]')];
        const labels = [];
        inputs.forEach(inp => { labels[Number(inp.dataset.outIdx)] = inp.value || ''; });
        node.data.outputs = labels;
        updateNodeHtml(node.id);
      });
    } catch {
      props.textContent = 'Select a node to edit.';
    }
  };

  editor.on('nodeSelected', (id) => renderProps(id));
  editor.on('nodeUnselected', () => {
    props.textContent = 'Select a node to edit.';
  });
}

function uuid() { return 'n-' + Math.random().toString(36).slice(2,9); }

// Palette: dynamic from templates
async function loadTemplates(){
  try{
    const r = await fetch('/api/templates');
    const j = await r.json();
    templates = j.templates || [];
    renderPalette();
  }catch(e){ console.warn('Failed to load templates', e); }
}

function renderPalette(){
  const pal = document.getElementById('palette');
  if(!pal) return;
  if(!templates || templates.length === 0){ pal.innerHTML = '<div class="text-xs text-gray-500">No templates found</div>'; return; }
  // Group by category
  const groups = {};
  templates.forEach(t => { const cat = t.category || 'other'; (groups[cat] ||= []).push(t); });
  const catOrder = Object.keys(groups).sort();
  pal.innerHTML = catOrder.map(cat => {
    const btns = groups[cat].map(t => `<button class="px-2 py-1 text-xs border rounded w-full text-left hover:bg-gray-50" onclick="addFromTemplate('${t.id.replace(/'/g, "&#39;")}')">
      <div class="font-medium">${t.name || t.id}</div>
      <div class="text-[11px] text-gray-500 truncate">${t.description || ''}</div>
    </button>`).join('');
    const title = cat[0].toUpperCase()+cat.slice(1);
    return `<div>
      <div class="text-xs uppercase tracking-wide text-gray-600 mb-1">${title}</div>
      <div class="flex flex-col gap-1">${btns}</div>
    </div>`;
  }).join('');
}

function nextNodePosition(){
  const baseX = 100; const baseY = 100;
  const count = Object.keys(editor.drawflow.drawflow.Home.data || {}).length;
  const x = baseX + (count % 5) * 160;
  const y = baseY + Math.floor(count / 5) * 120;
  return {x,y};
}

function addFromTemplate(tplId){
  const tpl = (templates || []).find(t => t.id === tplId);
  if(!tpl) return alert('Template not found: '+tplId);
  const {x,y} = nextNodePosition();
  const outCount = (tpl.initialOutputs && tpl.initialOutputs.length) ? tpl.initialOutputs.length : 0;
  const type = tpl.id;
  const data = {
    type,
    name: tpl.name || tpl.id,
    description: tpl.description || '',
    templateId: tpl.id,
    args: {...(tpl.input?.args || {})},
    outputs: [...(tpl.initialOutputs || [])]
  };
  const html = nodeHtml(type, data);
  const id = editor.addNode(type, 1, outCount, x, y, '', data, html, false);
  return id;
}

function nodeHtml(type, data){
  const title = data?.name || type;
  const outBadges = (data?.outputs || []).map(l => `<span class="px-1.5 py-0.5 bg-gray-200 rounded text-[10px] mr-1">${l}</span>`).join('');
  return `<div class="px-3 py-2 bg-gray-100 rounded-t border-b text-sm font-semibold">${title}</div>
  <div class="p-3 text-xs select-none">
    <div class="text-[11px] text-gray-500">${data?.description || ''}</div>
    ${outBadges?`<div class=\"mt-1\">${outBadges}</div>`:''}
  </div>`;
}

function updateNodeHtml(id){
  try{
    const node = editor.getNodeFromId(id);
    node.html = nodeHtml(node.name, node.data);
    editor.updateNodeDataFromId(id, node.data); // force re-render
  }catch{}
}

// Server interactions
async function checkHealth(){
  try{ const r = await fetch('/api/health'); const j = await r.json(); document.getElementById('health').textContent = j.status; }catch{ document.getElementById('health').textContent = 'offline'; }
}

async function listFlows(){
  try{ const r = await fetch('/api/flows'); const j = await r.json(); const ul = document.getElementById('flows'); ul.innerHTML = ''; j.forEach(id=>{ const li=document.createElement('li'); li.innerHTML=`<button class='underline' onclick=loadFlow('${id}')>${id}</button>`; ul.appendChild(li); }); }catch{}
}

async function loadFlow(id){
  const r = await fetch('/api/flows/' + id);
  const text = await r.text();
  try{ const data = jsyaml.load(text); const df = yamlToDrawflow(data); editor.import(df, false); document.getElementById('flowId').value = id; }catch(e){ alert('Failed to parse YAML: '+e); }
}

async function saveFlow(){
  if(!editor) return;
  const id = document.getElementById('flowId').value || 'untitled';
  const exported = editor.export();
  const yamlFlow = drawflowToYaml(exported, id);
  const yaml = jsyaml.dump(yamlFlow);
  const r = await fetch('/api/flows/'+id, { method:'POST', headers:{'Content-Type':'application/x-yaml'}, body: yaml });
  if(r.ok){ alert('Saved '+id); listFlows(); } else { alert('Save failed'); }
}

async function importFlow(file){
  const fd = new FormData(); fd.append('file', file);
  const r = await fetch('/api/flows/import', { method:'POST', body: fd });
  if(r.ok){ const j = await r.json(); listFlows(); alert('Imported '+j.imported); }
}

// Mapping helpers
function yamlToDrawflow(yamlObj){
  const f = yamlObj.flow || yamlObj || {};
  const nodes = f.nodes || [];
  const conns = f.connections || [];
  const data = {};

  // Index connections by from/to for port inference
  const fromMap = {};
  conns.forEach(c => {
    const k = c?.from?.node;
    if (!k) return;
    (fromMap[k] ||= []).push(c);
  });

  for(const n of nodes){
    const id = String(n.id);
    // Determine output labels
    let labels = [];
    const propsOuts = n?.properties?.outputs;
    if (Array.isArray(propsOuts)) {
      labels = propsOuts.map(String);
    }
    // Infer max port index used in connections: supports 'out', 'out1', 'out2', ...
    const used = fromMap[id] || [];
    let maxIdx = 0;
    used.forEach(c => {
      const p = c?.from?.port || 'out';
      const m = String(p).match(/^out(\d+)?$/);
      let idx = 1;
      if (m) { idx = m[1] ? parseInt(m[1],10) : 1; }
      if (idx > maxIdx) maxIdx = idx;
    });
    if (labels.length < maxIdx) {
      for (let i = labels.length; i < maxIdx; i++) labels[i] = `out${i+1}`;
    }

    // Build outputs object with N ports
    const outputs = {};
    const outCount = Math.max(labels.length, 0);
    if (outCount === 0) {
      // zero outputs allowed
    } else {
      for (let i = 1; i <= outCount; i++) {
        outputs[`output_${i}`] = { connections: [] };
      }
    }

    const node = {
      id,
      name: n.type || n.name || 'node',
      data: { type: n.type, name: n.properties?.name || n.type, description: n.properties?.description || '', outputs: labels },
      class: '',
      html: nodeHtml(n.type || 'node', { name: n.properties?.name || n.type, description: n.properties?.description || '', outputs: labels }),
      typenode: false,
      inputs: { input_1: { connections: [] } },
      outputs,
      pos_x: n.position?.x || 0,
      pos_y: n.position?.y || 0
    };
    data[id] = node;
  }

  // Apply connections mapping ports
  for(const c of conns){
    const from = c.from?.node; const to = c.to?.node;
    if(!from || !to) continue;
    if(!data[from] || !data[to]) continue;
    const port = c.from?.port || 'out';
    const m = String(port).match(/^out(\d+)?$/);
    let idx = 1; if (m) { idx = m[1] ? parseInt(m[1],10) : 1; }
    const outKey = `output_${idx}`;
    if (!data[from].outputs[outKey]) {
      data[from].outputs[outKey] = { connections: [] };
    }
    data[from].outputs[outKey].connections.push({ node: String(to), output: 'input_1' });
    data[to].inputs.input_1.connections.push({ node: String(from), input: outKey });
  }

  return { drawflow: { Home: { data } } };
}

function drawflowToYaml(dfExport, flowId){
  const data = (dfExport && dfExport.drawflow && dfExport.drawflow.Home && dfExport.drawflow.Home.data) || {};
  const nodes = [];
  const connections = [];

  for(const id of Object.keys(data)){
    const n = data[id];
    const outputs = (n?.data?.outputs && Array.isArray(n.data.outputs)) ? n.data.outputs : [];
    const properties = {
      name: n?.data?.name || n.name || '',
      description: n?.data?.description || '',
      args: n?.data?.args || {},
      outputs
    };
    nodes.push({
      id: String(n.id),
      type: n?.data?.type || n.name || 'node',
      position: { x: Number(n.pos_x) || 0, y: Number(n.pos_y) || 0 },
      properties
    });
  }

  for(const id of Object.keys(data)){
    const n = data[id];
    if(n.outputs){
      const outKeys = Object.keys(n.outputs).filter(k => k.startsWith('output_')).sort((a,b)=>{
        const ai = parseInt(a.split('_')[1],10)||0; const bi = parseInt(b.split('_')[1],10)||0; return ai-bi; });
      outKeys.forEach((outKey, idx) => {
        const outs = n.outputs[outKey].connections || [];
        const portName = (idx===0) ? 'out' : `out${idx+1}`;
        outs.forEach(o => {
          connections.push({
            from: { node: String(n.id), port: portName },
            to: { node: String(o.node), port: 'in' }
          });
        });
      });
    }
  }

  return { flow: { id: flowId, name: flowId, description: '', nodes, connections } };
}

// Initialize on DOM ready
window.addEventListener('DOMContentLoaded', () => {
  initDrawflow();
  checkHealth();
  listFlows();
  loadTemplates();
});

// Expose functions globally for HTML buttons
window.saveFlow = saveFlow;
window.importFlow = importFlow;
window.loadFlow = loadFlow;
window.listFlows = listFlows;
window.checkHealth = checkHealth;
window.addFromTemplate = addFromTemplate;
