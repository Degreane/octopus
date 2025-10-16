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
                        });
                        htmx.trigger(elt, 'htmx-set:success');
                    });
                }
            }
        }
    });
})();