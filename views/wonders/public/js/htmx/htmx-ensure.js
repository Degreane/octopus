(function () {
    htmx.defineExtension('htmx-ensure', {
        init: function (api) {
            var self = this;
            document.addEventListener('DOMContentLoaded', function () {
                self.processLoadEnsureElements();

                // Listen for DOM changes to handle dynamically added elements
                var observer = new MutationObserver(self.processLoadEnsureElements.bind(self));
                observer.observe(document.body, {childList: true, subtree: true});
                
                // Handle custom triggers
                var ensureObserver = new MutationObserver(function() {
                    document.querySelectorAll('[hx-ensure]').forEach(function(el) {
                        var triggerSpec = el.getAttribute("hx-ensure-trigger") || 'click';
                        if (triggerSpec !== 'click' && !el.hasAttribute('data-ensure-listener')) {
                            el.setAttribute('data-ensure-listener', 'true');
                            el.addEventListener(triggerSpec, function(evt) {
                                self.processElement('hx-ensure', el);
                            });
                        }
                    });
                });
                ensureObserver.observe(document.body, {childList: true, subtree: true, attributes: true, attributeFilter: ['hx-ensure']});
            });

            // Add event listeners to document for hx-ensure elements
            document.addEventListener('click', this.handleClick.bind(this));
        },

        processLoadEnsureElements: function () {
            document.querySelectorAll('[hx-ensure-load]').forEach(this.processElement.bind(this, 'hx-ensure-load'));
        },

        handleClick: function (event) {
            var element = event.target.closest('[hx-ensure]');
            if (element) {
                var triggerSpec = element.getAttribute("hx-ensure-trigger") || 'click';
                if (event.type === triggerSpec) {
                    this.processElement('hx-ensure', element);
                }
            }
        },

        processElement: function (attrName, elt) {
            var ensureAttr = elt.getAttribute(attrName);
            if (ensureAttr) {
                var instructions = ensureAttr.split(',');
                instructions.forEach(function (instruction) {
                    var colonIndex = instruction.indexOf(':');
                    if (colonIndex === -1) return;

                    var selector = instruction.substring(0, colonIndex).trim();
                    var actions = instruction.substring(colonIndex + 1).trim();
                    var target = document.querySelector(selector);

                    if (target) {
                        actions.split(';').forEach(function (action) {
                            var act = action.trim();

                            // Check if it's an attribute operation (contains '=')
                            if (act.includes('=')) {
                                var attrMatch = act.match(/^(!?)([^=]+)=(.*)$/);
                                if (attrMatch) {
                                    var removeAttr = attrMatch[1] === '!';
                                    var attrName = attrMatch[2].trim();
                                    var attrValue = attrMatch[3].trim();

                                    if (removeAttr) {
                                        target.removeAttribute(attrName);
                                    } else {
                                        target.setAttribute(attrName, attrValue);
                                    }
                                }
                                return;
                            }

                            // Support formats for classes:
                            //  - "className"          => add immediately
                            //  - "!className"         => remove immediately
                            //  - "className 3s"       => add after 3s
                            //  - "!className 1.5s"    => remove after 1.5s
                            var timedMatch = act.match(/^(!?)(\S+)\s+(\d+(?:\.\d+)?)s$/);

                            if (timedMatch) {
                                var removeTimed = timedMatch[1] === '!';
                                var classNameTimed = timedMatch[2];
                                var seconds = parseFloat(timedMatch[3]);

                                setTimeout(function () {
                                    if (removeTimed) {
                                        target.classList.remove(classNameTimed);
                                    } else {
                                        target.classList.add(classNameTimed);
                                    }
                                }, seconds * 1000);
                            } else {
                                var remove = act.startsWith('!');
                                var className = remove ? act.substring(1) : act;

                                if (remove) {
                                    target.classList.remove(className);
                                } else {
                                    target.classList.add(className);
                                }
                            }
                        });
                    }
                });
            }
        },

        onEvent: function (name, evt) {
            if (name === 'htmx:afterSettle') {
                this.processLoadEnsureElements();
            }
        }
    });

})();