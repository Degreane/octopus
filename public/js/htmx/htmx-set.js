(function () {
    htmx.defineExtension('htmx-set', {
        onEvent: function (name, evt) {
            if (name === "htmx:afterProcessNode") {
                var elt = evt.detail.elt;
                var setSpec = elt.getAttribute("hx-set");

                if (setSpec) {
                    elt.addEventListener('click', function () {
                        var specs = setSpec.split(/,,/);
                        specs.forEach(function (spec) {
                            var [id, value] = spec.trim().split(':');
                            var element = document.querySelector(id);
                            if (element) {
                                // New protocol:
                                // hx-set="#id:get;;value;;target;;swap;;trigger"
                                // -> sets hx-get, hx-target, hx-swap, hx-trigger on the target element
                                if (typeof value === 'string' && value.indexOf('get;;') === 0) {
                                    var parts = value.split(';;');
                                    // parts[0] = "get", parts[1]=value, parts[2]=target, parts[3]=swap, parts[4]=trigger
                                    var getVal = parts[1] || '';
                                    var targetVal = parts[2] || '';
                                    var swapVal = parts[3] || '';
                                    var triggerVal = parts[4] || '';
                                    element.setAttribute('hx-get', getVal);
                                    if (targetVal) element.setAttribute('hx-target', targetVal);
                                    if (swapVal) element.setAttribute('hx-swap', swapVal);
                                    if (triggerVal) element.setAttribute('hx-trigger', triggerVal);

                                    // Ensure htmx processes the element after dynamic trigger assignment
                                    // so new hx-* attributes are recognized and bound.
                                    try {
                                        if (typeof htmx.process === 'function') {
                                            htmx.process(element);
                                        } else {
                                            // Fallback hint (rarely needed)
                                            htmx.trigger(document.body, 'htmx:load', { elt: element });
                                        }
                                    } catch (e) {
                                        // no-op
                                    }
                                } else {
                                    if (element.tagName.toLowerCase() === 'img') {
                                        element.src = value;
                                    } else if (element.tagName.toLowerCase() === 'input') {
                                        element.value = value;
                                    } else if (element.tagName.toLowerCase() === 'a') {
                                        if (value.startsWith('loc;;')) {
                                            var value2 = value.split(/;;/);
                                            document.location = value2[1];
                                        } else if (value.startsWith('href;;')) {
                                            var value2 = value.split(/;;/);
                                            element.href = value2[1];
                                        } else {
                                            element.textContent = value;
                                        }
                                    } else {
                                        element.textContent = value;
                                    }
                                }
                            }
                        });
                        htmx.trigger(elt, 'htmx-set:success');
                    });
                }
            }
        }
    });
})();