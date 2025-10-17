htmx.defineExtension('htmx-copy', {
    onEvent: function (name, evt) {
        if (name === "htmx:afterProcessNode") {
            var target = evt.detail.elt;
            if (target.hasAttribute('hx-copy')) {
                var copyValue = target.getAttribute('hx-copy');
                target.addEventListener('click', function(e) {
                    var textToCopy;

                    if (copyValue.startsWith('#')) {
                        // Copy from another element
                        var sourceElement = document.querySelector(copyValue);
                        if (sourceElement) {
                            if (sourceElement.tagName === 'INPUT' || sourceElement.tagName === 'TEXTAREA') {
                                textToCopy = sourceElement.value;
                            } else {
                                textToCopy = sourceElement.textContent || sourceElement.innerText;
                            }
                        }
                    } else {
                        // Copy literal text
                        textToCopy = copyValue;
                    }

                    if (textToCopy) {
                        if (navigator.clipboard && navigator.clipboard.writeText) {
                            navigator.clipboard.writeText(textToCopy).then(function() {
                                console.log('Text copied to clipboard');
                            }).catch(function(err) {
                                console.error('Failed to copy text: ', err);
                            });
                        } else {
                            // Fallback for older browsers
                            var textArea = document.createElement("textarea");
                            textArea.value = textToCopy;
                            document.body.appendChild(textArea);
                            textArea.focus();
                            textArea.select();
                            try {
                                document.execCommand('copy');
                                console.log('Text copied to clipboard (fallback)');
                            } catch (err) {
                                console.error('Fallback copy failed: ', err);
                            }
                            document.body.removeChild(textArea);
                        }
                    }
                });
            }

            // Add paste functionality
            if (target.hasAttribute('hx-paste')) {
                var pasteTarget = target.getAttribute('hx-paste');
                target.addEventListener('click', function(e) {
                    if (pasteTarget.startsWith('#')) {
                        var targetElement = document.querySelector(pasteTarget);
                        if (targetElement) {
                            if (navigator.clipboard && navigator.clipboard.readText) {
                                navigator.clipboard.readText().then(function(clipboardText) {
                                    if (targetElement.tagName === 'INPUT' || targetElement.tagName === 'TEXTAREA') {
                                        targetElement.value = clipboardText;
                                    } else {
                                        targetElement.textContent = clipboardText;
                                    }
                                    targetElement.focus();
                                    console.log('Text pasted from clipboard');
                                }).catch(function(err) {
                                    console.error('Failed to read clipboard: ', err);
                                });
                            } else {
                                // Fallback for older browsers
                                targetElement.focus();
                                try {
                                    document.execCommand('paste');
                                    console.log('Text pasted from clipboard (fallback)');
                                } catch (err) {
                                    console.error('Fallback paste failed: ', err);
                                }
                            }
                        }
                    }
                });
            }
        }
    }
});