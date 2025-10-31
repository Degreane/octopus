### df-app.js — Detailed Documentation

#### File purpose
This file provides a thin wrapper around the Drawflow library to power the visual flow editor. It:
- Initializes the Drawflow editor instance
- Adds nodes via palette buttons
- Loads/saves flows to the backend in the YAML format your API expects
- Converts between your YAML schema and Drawflow’s internal JSON format
- Updates a simple properties panel
- Exposes a few functions globally so the HTML buttons can call them

---

#### Functions overview
- `initDrawflow()`: Bootstraps the Drawflow editor, sets options, binds selection events to update the properties panel.
- `uuid()`: Generates a simple client-side ID string (not used by Drawflow when `useuuid=true`, but kept as a utility).
- `addNode(name, type)`: Adds a new node with one input and one output to the canvas at a computed position.
- `checkHealth()`: Calls the backend health endpoint and updates UI.
- `listFlows()`: Lists available flows from the backend and renders clickable items to load them.
- `loadFlow(id)`: Fetches a YAML flow by ID, converts it to Drawflow JSON, and imports it into the editor.
- `saveFlow()`: Exports Drawflow JSON, converts to your YAML schema, and saves to backend under the given ID.
- `importFlow(file)`: Uploads a YAML file to the backend import endpoint; the backend saves it to the flows folder.
- `yamlToDrawflow(yamlObj)`: Maps your YAML schema to Drawflow’s JSON structure (one input, one output per node convention).
- `drawflowToYaml(dfExport, flowId)`: Maps Drawflow JSON back to your YAML schema for persistence.
- DOMContentLoaded handler: Initializes editor and auxiliary UI on page ready.
- Global bindings: Expose core functions under `window.*` so the HTML can call them.

---

#### Line-by-line documentation

1-2. File header comments describing the purpose of the wrapper and that it replaces a previous custom implementation.

4. `let editor = null;`
   - Declares a module-level variable that will hold the Drawflow editor instance after initialization.

6. `function initDrawflow() {`
   - Entry point to set up Drawflow and bind UI events.

7. `const el = document.getElementById('drawflow');`
   - Looks up the canvas container DIV in `index.html` where Drawflow will render.

8. `editor = new Drawflow(el);`
   - Instantiates the Drawflow editor bound to that container element.

9-13. Editor basic configuration.
   - `10: editor.start();` Boots the editor and attaches all necessary event handlers and DOM structure.
   - `13: editor.useuuid = true;` Instructs Drawflow to use UUIDs for node IDs, giving stable string IDs.

15. Comment: the following block wires a simple properties panel to node selection.

16. `const props = document.getElementById('props');`
   - Grabs the properties panel element to display current node info.

17. `editor.on('nodeSelected', (id) => {`
   - Subscribes to Drawflow’s `nodeSelected` event; receives the selected node id.

18-29. Selection handler body.
   - `19: const node = editor.getNodeFromId(id);` Fetch the full node record from Drawflow.
   - `20: const type = node?.data?.type || node?.name || '';` Prefer the node’s `data.type`, fallback to `name`.
   - `21-25:` Populate the properties panel with node id and type using a small HTML fragment.
   - `26-28:` If anything fails, show a default "Select a node" message (e.g., if node not found).

30-32. `editor.on('nodeUnselected', ...)`
   - When a node is unselected, update the properties panel back to the default message.

33. End of `initDrawflow()`.

35. `function uuid() { return 'n-' + Math.random().toString(36).slice(2,9); }`
   - Simple random ID generator for non-critical IDs. Kept as a utility; Drawflow’s `useuuid` handles node IDs.

37. Comment: Palette function for adding nodes.

38. `function addNode(name, type) {`
   - Public function called by palette buttons to add a new node of a given `type` to the editor.

39. `if (!editor) return;`
   - Guard if the editor isn’t initialized yet.

40. `const baseX = 100; const baseY = 100;`
   - Establishes a starting location for newly added nodes.

42. `const count = Object.keys(editor.drawflow.drawflow.Home.data || {}).length;`
   - Counts existing nodes to space out new nodes.

43-44. Compute `x` and `y` coordinates for the new node in a simple grid layout.

46-47. `html` template string used as the node’s body content (title + helper text).

49. Comment: clarify that Drawflow’s `name` parameter is the node type we show on the title.

50. `const id = editor.addNode(type, 1, 1, x, y, '', { type, name }, html, false);`
   - Adds a node to Drawflow:
     - `type` as display name,
     - 1 input, 1 output,
     - positioned at `(x, y)`,
     - CSS class is an empty string,
     - `data` payload: `{ type, name }`,
     - `html` is the node’s content,
     - `typenode=false` indicates standard HTML content, not a registered component.

51-52. Return the new node id and close the function.

54. Comment: Server interaction helpers start here.

55-57. `async function checkHealth(){ ... }`
   - Calls `/api/health`, parses JSON, updates `#health` display; falls back to `offline` if request fails.

59-61. `async function listFlows(){ ... }`
   - Calls `/api/flows` to get an array of flow IDs.
   - Clears the `#flows` list and appends a button per id that calls `loadFlow(id)` on click.
   - Silently ignores errors (e.g., server down).

63. `async function loadFlow(id){`
   - Fetches a specific flow by id from `/api/flows/:id` (YAML response).

64-66. Fetch YAML as text, parse with js-yaml, map to Drawflow JSON via `yamlToDrawflow`, then `editor.import(df, false)`:
   - `false` avoids firing Drawflow’s `import` event. Set `#flowId` input value to the loaded id.
   - If parsing fails, alerts the user.

69. `async function saveFlow(){`
   - Exports the current Drawflow graph, converts to YAML schema, and posts to the backend.

70-77. Save body.
   - `71:` Determine the flow id from `#flowId` (defaults to `untitled`).
   - `72:` `editor.export()` returns Drawflow JSON.
   - `73-74:` Convert to YAML-compatible structure via `drawflowToYaml` and serialize with `jsyaml.dump`.
   - `75:` POST YAML to `/api/flows/:id` with `Content-Type: application/x-yaml`.
   - `76:` On success, alert and refresh the flow list; else, alert failure.

79. `async function importFlow(file){`
   - Upload a YAML file to the backend importer.

80-83. Build a `FormData`, POST to `/api/flows/import`, and on success report the saved name and refresh the list.

85. Comment: Mapping helpers start here.

86. `function yamlToDrawflow(yamlObj){`
   - Converts your YAML schema to Drawflow’s structure with a single `Home` module.

87-91. Extract `flow` content (or accept a top-level object), default arrays, prepare the `data` object that will hold nodes keyed by id.

93-107. Build Drawflow node records for each YAML node.
   - `94-106:` For each YAML node `n`, create an entry keyed by `n.id` with:
     - `id, name, data, class, html`: the node’s metadata and visible content
     - `inputs: { input_1 }` and `outputs: { output_1 }` by convention
     - `pos_x/pos_y`: mapped from `position.x / position.y` (defaults to 0)

109-118. Add connections from YAML to Drawflow.
   - Iterate `conns` and translate `from: out -> to: in` onto `outputs.output_1.connections` and `inputs.input_1.connections`.
   - Skip invalid references.

120-121. Return the Drawflow export object `{ drawflow: { Home: { data } } }` ready for `editor.import()`.

123. `function drawflowToYaml(dfExport, flowId){`
   - Maps back from Drawflow export JSON to your YAML schema.

124. Extract the `Home.data` dictionary from the export.

125-136. Nodes conversion.
   - Iterate keys of `data` and push a YAML `node` with `id, type (prefer data.type), position, properties.name`.
   - Coerce numeric positions; default to 0.

138-149. Connections conversion.
   - For each node, inspect `outputs.output_1.connections` and add a YAML connection:
     - `from: { node: n.id, port: 'out' }`
     - `to:   { node: o.node, port: 'in' }`

151-152. Return the YAML object `{ flow: { id, name, description, nodes, connections } }` expected by the backend.

154-159. DOM ready handler.
   - On `DOMContentLoaded`, initialize Drawflow (`initDrawflow`), check server health, and list flows.

161-167. Global bindings.
   - Attach key functions to `window.*` so inline HTML event handlers in `index.html` can call them.

168. (intentionally blank at EOF)

---

#### Additional notes and edge cases
- Node IDs: With `editor.useuuid = true`, Drawflow assigns stable UUIDs for new nodes. When loading from YAML, `yamlToDrawflow` preserves the IDs from YAML (keys in `data`), so round-tripping keeps IDs intact.
- Single input/output convention: The mapping assumes one input (`input_1`) and one output (`output_1`) per node. If you later add multi-port nodes, you’ll need to extend both mapping helpers to enumerate ports and encode port names in the YAML `port` fields accordingly.
- Import semantics: `editor.import(df, false)` replaces the current graph with the provided data without emitting an `import` event. If you want to react to imports, set the second parameter to `true`.
- Error handling: Network and parsing errors in `loadFlow`, `saveFlow`, and `listFlows` are handled minimally via `alert` or silent failure. Consider enhancing UX with inline notifications.
- Security: This is a client-side tool used in development contexts; production deployments may want stricter CORS and CSRF protections on the backend.


---

### Templates system

The Node Palette now loads node templates from a YAML file on the server.

- Backend endpoint: `GET /api/templates` returns a JSON object `{ templates: Template[] }` parsed from `flows/templates.yaml`.
- File location: `flows/templates.yaml` with the following schema:

```
templates:
  - id: string                 # template ID, used as node `type`
    name: string               # human-readable template name
    description: string        # shown in palette and properties
    category: string           # one of: control, functional, config, display, ...
    input:
      variables: boolean       # whether inputs accept variables (UI hint)
      args:                    # default arguments for the node
        key: value
    initialOutputs:            # zero or more output port labels
      - outLabel1
      - outLabel2
```

- Frontend behavior:
  - On load, `loadTemplates()` fetches the endpoint and renders buttons grouped by `category` in the Palette.
  - Clicking a template calls `addFromTemplate(template.id)` which creates a node with 1 input and N outputs based on `initialOutputs`.
  - Node `data` contains: `{ type, name, description, templateId, args, outputs }`.
  - The Properties panel lets users:
    - Edit arbitrary key/value arguments (add/remove/save).
    - Add/remove output ports dynamically (updates both the visual ports and stored labels), and rename output labels.

- YAML import/export changes:
  - Drawflow export now persists extra fields under `node.properties`: `name`, `description`, `args`, and `outputs` (array of output labels, aligned by index with ports `out`, `out2`, `out3`, ...).
  - Connections `from.port` is exported as `out` for the first output, and `out2`, `out3`, ... for subsequent outputs.
  - YAML → Drawflow conversion reconstructs the correct number of output ports by reading `properties.outputs` or inferring from connections.

- Backward compatibility:
  - Older flows that only used a single output (`out`) still load correctly; additional outputs are inferred as needed from connections.

- Example `flows/templates.yaml` shipped with the project includes categories:
  - control: `control.if`, `control.loop`
  - functional: `func.define`
  - config: `config.read`
  - display: `display.log`

Usage tips:
- You can customize `flows/templates.yaml` to add your own node types and categories; no frontend code changes required.
- If the templates file is missing, the Palette will show "No templates found".
