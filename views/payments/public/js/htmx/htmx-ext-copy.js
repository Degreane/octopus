(function () {
    htmx.defineExtension('htmx-copy', {
        onEvent: function (name, evt) {
            if (name === "htmx:afterProcessNode") {
                var elt = evt.detail.elt;
                var copySpec = elt.getAttribute("hx-copy");

                if (copySpec) {
                    elt.addEventListener('click', function () {
                        var [targetSelector, duration, message] = copySpec.split(':');
                        // var targetSelector = copySpec.trim();

                        if (targetSelector.charAt(0) !== '#') {
                            // console.error("HTMX ext-copy: unsupported selector '" + targetSelector + "'. Only ID selectors starting with '#' are supported.");
                            return;
                        }

                        var targetElement = document.querySelector(targetSelector);

                        if (targetElement) {
                            if (targetElement.type == 'textarea') {
                                var textToCopy = targetElement.value;
                            } else {
                                var textToCopy = targetElement.innerText;
                            }
                            // var textToCopy = targetElement.innerText;

                            navigator.clipboard.writeText(textToCopy).then(function () {
                                //   console.log("Content copied to clipboard");
                                htmx.trigger(elt, 'clipboard:success');
                                if (duration && message) {
                                    var infoMessage = document.createElement('div');
                                    infoMessage.className = 'fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-blue-500 text-white px-4 py-2 rounded shadow-lg opacity-0 transition-opacity duration-600';
                                    infoMessage.innerHTML = message;
                                    document.body.appendChild(infoMessage);
                                    // Fade in
                                    setTimeout(() => infoMessage.classList.add('opacity-100'), 10);
                                    // Fade out and remove
                                    setTimeout(function () {
                                        infoMessage.classList.remove('opacity-100');
                                        setTimeout(() => document.body.removeChild(infoMessage), 600);
                                    }, (parseFloat(duration) * 1000) - 600);
                                }
                            }).catch(function (err) {
                                //   console.error("Failed to copy content: ", err);
                                htmx.trigger(elt, 'clipboard:error');
                            });
                        } else {
                            // console.warn("HTMX ext-copy: selector '" + targetSelector + "' not found in document.");
                            htmx.trigger(elt, 'clipboard:error');
                        }
                    });
                }
            }
        }
    });
})();
