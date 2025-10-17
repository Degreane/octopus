(function () {
    htmx.defineExtension('htmx-set', {
        onEvent: function (name, evt) {
            if (name === "htmx:afterProcessNode") {
                var elt = evt.detail.elt;
                var setSpec = elt.getAttribute("hx-set");

                if (setSpec) {
                    var triggerSpec = elt.getAttribute("hx-set-trigger") || 'click';
                    elt.addEventListener(triggerSpec, function () {
                        var specs = setSpec.split(/,,/);
                        specs.forEach(function (spec) {
                            var [id, value] = spec.trim().split(':');
                            var element = document.querySelector(id);
                            if (element) {
                                // Handle attribute setting with att;; prefix
                                if (value.startsWith('att;;')) {
                                    var attrSpec = value.substring(5); // Remove 'att;;'
                                    // Parse attribute=value format
                                    var attrMatch = attrSpec.match(/^([^=]+)=(.+)$/);
                                    if (attrMatch) {
                                        var attrName = attrMatch[1].trim();
                                        var attrValue = attrMatch[2].trim();
                                        // Remove surrounding quotes if present
                                        attrValue = attrValue.replace(/^['"]|['"]$/g, '');

                                        // Handle toggle for boolean-like attributes
                                        if (attrValue === 'toggle') {
                                            var currentValue = element.getAttribute(attrName);
                                            if (currentValue === 'true') {
                                                attrValue = 'false';
                                            } else {
                                                attrValue = 'true';
                                            }
                                        }

                                        element.setAttribute(attrName, attrValue);
                                    }
                                } else if (element.tagName.toLowerCase() === 'img') {
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
                        });
                        htmx.trigger(elt, 'htmx-set:success');
                    });
                }
            }
        }
    });
})();