let state = {
  nodes: [],
  connections: [],
  selected: null,
  pan: {active:false, startX:0, startY:0, scrollLeft:0, scrollTop:0},
  connectDrag: {active:false, from:null, mouse:{x:0,y:0}, hover:null}
};

function uuid(){ return 'n-' + Math.random().toString(36).slice(2,9); }

function addNode(name, type){
  const id = uuid();
  const node = { id, type, name, position:{x: 100 + state.nodes.length*40, y: 100}, properties:{} };
  state.nodes.push(node);
  render();
}

function render(){
  const nodesEl = document.getElementById('nodes');
  nodesEl.innerHTML = '';
  state.nodes.forEach(n => {
    const el = document.createElement('div');
    el.className = 'node absolute bg-white rounded border min-w-[160px] cursor-move';
    el.style.left = n.position.x + 'px';
    el.style.top = n.position.y + 'px';
    el.dataset.id = n.id;
    el.onmousedown = (e)=> startDrag(e, n.id);
    el.innerHTML = `
      <div class="px-3 py-2 bg-gray-100 rounded-t border-b text-sm font-semibold">${n.type}</div>
      <div class="p-3 text-xs select-none">
        <div class="flex items-center gap-2 mb-2">
          <span class="port in string" data-node="${n.id}" data-port="in" title="input"></span>
          <span>in</span>
          <span class="ml-auto port out number" data-node="${n.id}" data-port="out" title="output"></span>
        </div>
        <div>ID: ${n.id}</div>
      </div>`;
    el.onclick = ()=> selectNode(n.id);
    nodesEl.appendChild(el);
  });
  // Attach port events
  setTimeout(()=>{
    document.querySelectorAll('#nodes .port.out').forEach(p=>{
      p.onmousedown = (e)=> startConnectionDrag(e, p.dataset.node, p.dataset.port);
      p.onpointerdown = (e)=> startConnectionDrag(e, p.dataset.node, p.dataset.port);
    });
    document.querySelectorAll('#nodes .port.in').forEach(p=>{
      // Prevent starting a node drag when clicking on an input port
      p.onmousedown = (e)=> { e.stopPropagation(); };
      p.onpointerdown = (e)=> { e.stopPropagation(); };
      // Track hover target robustly during connection drag
      p.onpointerenter = (e)=> { if(state.connectDrag.active){ state.connectDrag.hover = e.currentTarget; e.currentTarget.classList.add('targetable'); } };
      p.onpointerleave = (e)=> { if(state.connectDrag.active && state.connectDrag.hover===e.currentTarget){ state.connectDrag.hover = null; e.currentTarget.classList.remove('targetable'); } };
      // If we release the mouse over an input while connecting, finalize the connection
      p.onmouseup = (e)=> {
        e.stopPropagation();
        if(state.connectDrag && state.connectDrag.active){
          finishConnectionToTarget(e.currentTarget);
        }
      };
      p.onpointerup = (e)=> {
        e.stopPropagation();
        if(state.connectDrag && state.connectDrag.active){
          finishConnectionToTarget(e.currentTarget);
        }
      };
    });
  },0);
  renderProps();
  renderWires();
}

function selectNode(id){
  state.selected = state.nodes.find(n=>n.id===id) || null;
  renderProps();
}

function renderProps(){
  const p = document.getElementById('props');
  const n = state.selected;
  if(!n){ p.textContent = 'Select a node to edit.'; return; }
  p.innerHTML = `
    <div class="space-y-2">
      <div class="text-sm">Node: <span class="font-mono">${n.id}</span></div>
      <label class="block text-xs">Name
        <input class="border w-full px-2 py-1 text-sm" value="${n.name||''}" oninput="n.properties.name=this.value" />
      </label>
    </div>`;
}

// Drag handling
let drag = {active:false, id:null, dx:0, dy:0, startX:0, startY:0};
function startDrag(e, id){
  drag.active = true; drag.id = id; drag.startX = e.clientX; drag.startY = e.clientY;
  const n = state.nodes.find(n=>n.id===id);
  drag.dx = n.position.x; drag.dy = n.position.y;
  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', endDrag);
}
function onDrag(e){
  if(!drag.active) return;
  const n = state.nodes.find(n=>n.id===drag.id);
  n.position.x = drag.dx + (e.clientX - drag.startX);
  n.position.y = drag.dy + (e.clientY - drag.startY);
  render();
}
function endDrag(){
  drag.active = false; drag.id = null;
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', endDrag);
}

// Pan handling
function startPan(e){ if(e.target.id!=='canvas') return; state.pan={active:true,startX:e.clientX,startY:e.clientY,scrollLeft:e.target.scrollLeft,scrollTop:e.target.scrollTop}; }
function onPan(e){ const c = document.getElementById('canvas'); if(!state.pan.active) return; c.scrollLeft = state.pan.scrollLeft - (e.clientX - state.pan.startX); c.scrollTop = state.pan.scrollTop - (e.clientY - state.pan.startY); }
function endPan(){ state.pan.active=false; }

// Geometry helpers
function svgCoordsFromClient(clientX, clientY){
  const svg = document.getElementById('wires');
  const r = svg.getBoundingClientRect();
  return { x: clientX - r.left, y: clientY - r.top };
}
function getPortCenter(nodeId, port){
  const el = document.querySelector(`#nodes .port[data-node="${nodeId}"][data-port="${port}"]`);
  const svg = document.getElementById('wires');
  if(!el || !svg){ return null; }
  const rb = el.getBoundingClientRect();
  const rs = svg.getBoundingClientRect();
  return { x: rb.left - rs.left + rb.width/2, y: rb.top - rs.top + rb.height/2 };
}
// Drawflow-inspired curvature for cubic Bezier path between p1 (output) and p2 (input)
function curvaturePath(p1, p2, curvature = 0.5){
  const dx = Math.abs(p2.x - p1.x);
  const direction = p1.x <= p2.x ? 1 : -1; // open/close behavior
  const cx1 = p1.x + dx * curvature * direction;
  const cy1 = p1.y;
  const cx2 = p2.x - dx * curvature * direction;
  const cy2 = p2.y;
  return `M ${p1.x} ${p1.y} C ${cx1} ${cy1} ${cx2} ${cy2} ${p2.x} ${p2.y}`;
}

// Connection drag handlers
function startConnectionDrag(e, nodeId, port){
  if(port !== 'out') return;
  e.stopPropagation();
  e.preventDefault();
  state.connectDrag.active = true;
  state.connectDrag.from = { node: nodeId, port: 'out' };
  const p = svgCoordsFromClient(e.clientX, e.clientY);
  state.connectDrag.mouse = p;
  document.addEventListener('mousemove', onConnectionMove);
  document.addEventListener('mouseup', endConnectionDrag);
  document.addEventListener('pointermove', onConnectionMove);
  document.addEventListener('pointerup', endConnectionDrag);
  // highlight input ports
  document.querySelectorAll('#nodes .port.in').forEach(el=> el.classList.add('targetable'));
  renderWires();
}
function onConnectionMove(e){
  if(!state.connectDrag.active) return;
  state.connectDrag.mouse = svgCoordsFromClient(e.clientX, e.clientY);
  renderWires();
}
function endConnectionDrag(e){
  if(!state.connectDrag.active) return;
  // Prefer tracked hover input if available
  let target = state.connectDrag.hover;
  if(!target && e && typeof e.clientX === 'number'){
    const els = document.elementsFromPoint(e.clientX, e.clientY);
    target = els.find(el => el.classList && el.classList.contains('port') && el.classList.contains('in')) || null;
  }
  if(target){
    finishConnectionToTarget(target);
  }
  document.removeEventListener('mousemove', onConnectionMove);
  document.removeEventListener('mouseup', endConnectionDrag);
  document.removeEventListener('pointermove', onConnectionMove);
  document.removeEventListener('pointerup', endConnectionDrag);
  document.querySelectorAll('#nodes .port.in').forEach(el=> el.classList.remove('targetable'));
  state.connectDrag.active = false; state.connectDrag.from = null; state.connectDrag.hover = null;
  renderWires();
}

function finishConnectionToTarget(target){
  const toNode = target && target.dataset ? target.dataset.node : null;
  if(!toNode) return;
  if(!state.connectDrag || !state.connectDrag.from) return;
  if(toNode === state.connectDrag.from.node) return; // prevent self-connection
  const conn = { from: { node: state.connectDrag.from.node, port: 'out' }, to: { node: toNode, port: 'in' } };
  const dup = state.connections.some(c=> c.from.node===conn.from.node && c.from.port===conn.from.port && c.to.node===conn.to.node && c.to.port===conn.to.port);
  if(!dup){ state.connections.push(conn); }
}

// Wire rendering (existing + preview)
function renderWires(){
  const svg = document.getElementById('wires');
  svg.innerHTML = '';
  // existing connections
  state.connections.forEach((c)=>{
    const p1 = getPortCenter(c.from.node, c.from.port);
    const p2 = getPortCenter(c.to.node, c.to.port);
    if(p1 && p2){
      const path = document.createElementNS('http://www.w3.org/2000/svg','path');
      path.setAttribute('d', curvaturePath(p1,p2));
      path.setAttribute('class','wire');
      svg.appendChild(path);
    }
  });
  // preview
  if(state.connectDrag.active && state.connectDrag.from){
    const p1 = getPortCenter(state.connectDrag.from.node, 'out');
    const p2 = state.connectDrag.mouse;
    if(p1 && p2){
      const path = document.createElementNS('http://www.w3.org/2000/svg','path');
      path.setAttribute('d', curvaturePath(p1,p2));
      path.setAttribute('class','wire active');
      svg.appendChild(path);
    }
  }
}

window.addEventListener('resize', ()=> renderWires());

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
  try{ const data = jsyaml.load(text); applyFlow(data); document.getElementById('flowId').value = id; }catch(e){ alert('Failed to parse YAML: '+e); }
}
function applyFlow(data){
  const f = data.flow || data; // allow either structure for resilience
  state.nodes = (f.nodes||[]).map(n=>({ id:n.id, type:n.type, name:n.type, position:n.position||{x:0,y:0}, properties:n.properties||{} }));
  state.connections = f.connections||[];
  render();
}
async function saveFlow(){
  const id = document.getElementById('flowId').value || 'untitled';
  const flow = { flow: { id, name: id, description: '', nodes: state.nodes.map(n=>({id:n.id,type:n.type,position:n.position,properties:n.properties})), connections: state.connections } };
  const yaml = jsyaml.dump(flow);
  const r = await fetch('/api/flows/'+id, { method:'POST', headers:{'Content-Type':'application/x-yaml'}, body: yaml });
  if(r.ok){ alert('Saved '+id); listFlows(); } else { alert('Save failed'); }
}
async function importFlow(file){
  const fd = new FormData(); fd.append('file', file);
  const r = await fetch('/api/flows/import', { method:'POST', body: fd });
  if(r.ok){ const j = await r.json(); listFlows(); alert('Imported '+j.imported); }
}

// Initialize after DOM is ready
window.addEventListener('DOMContentLoaded', () => {
  checkHealth(); listFlows(); render();
  const canvas = document.getElementById('canvas');
  if(canvas){ canvas.addEventListener('scroll', ()=> renderWires(), { passive: true }); }
});
