(function () {
    htmx.defineExtension('htmx-content', {
        init: function (api) {
            var self = this;
            document.addEventListener('DOMContentLoaded', function () {
                self.processLoadContentElements();

                // Listen for DOM changes to handle dynamically added elements
                var observer = new MutationObserver(self.processLoadContentElements.bind(self));
                observer.observe(document.body, {childList: true, subtree: true});
                
                // Handle custom triggers
                var contentObserver = new MutationObserver(function() {
                    var selfRef = this; // Save reference to the extension object
                    document.querySelectorAll('[hx-content]').forEach(function(el) {
                        var triggerSpec = el.getAttribute("hx-content-trigger") || 'load';
                        if (triggerSpec !== 'load' && !el.hasAttribute('data-content-listener')) {
                            el.setAttribute('data-content-listener', 'true');
                            el.addEventListener(triggerSpec, function(evt) {
                                var delay = parseFloat(el.getAttribute("hx-content-delay")) || 0;
                                if (delay > 0) {
                                    setTimeout(function() {
                                        selfRef.processElement(el);
                                    }, delay * 1000);
                                } else {
                                    selfRef.processElement(el);
                                }
                            });
                        }
                    });
                }.bind(this));
                contentObserver.observe(document.body, {childList: true, subtree: true, attributes: true, attributeFilter: ['hx-content']});
            });

            // Add event listeners to document for hx-content elements
            document.addEventListener('click', this.handleClick.bind(this));
            document.addEventListener('mouseover', this.handleMouseOver.bind(this));
            document.addEventListener('change', this.handleChange.bind(this));
        },

        processLoadContentElements: function () {
            document.querySelectorAll('[hx-content]').forEach(function(el) {
                var triggerSpec = el.getAttribute("hx-content-trigger") || 'load';
                if (triggerSpec === 'load') {
                    var delay = parseFloat(el.getAttribute("hx-content-delay")) || 0;
                    if (delay > 0) {
                        setTimeout(function() {
                            this.processElement(el);
                        }.bind(this), delay * 1000);
                    } else {
                        this.processElement(el);
                    }
                }
            }.bind(this));
        },

        handleClick: function (event) {
            var element = event.target.closest('[hx-content]');
            if (element) {
                var triggerSpec = element.getAttribute("hx-content-trigger") || 'load';
                if (triggerSpec === 'click') {
                    var delay = parseFloat(element.getAttribute("hx-content-delay")) || 0;
                    if (delay > 0) {
                        setTimeout(function() {
                            this.processElement(element);
                        }.bind(this), delay * 1000);
                    } else {
                        this.processElement(element);
                    }
                }
            }
        },

        handleMouseOver: function (event) {
            var element = event.target.closest('[hx-content]');
            if (element) {
                var triggerSpec = element.getAttribute("hx-content-trigger") || 'load';
                if (triggerSpec === 'hover') {
                    var delay = parseFloat(element.getAttribute("hx-content-delay")) || 0;
                    if (delay > 0) {
                        setTimeout(function() {
                            this.processElement(element);
                        }.bind(this), delay * 1000);
                    } else {
                        this.processElement(element);
                    }
                }
            }
        },

        handleChange: function (event) {
            var element = event.target.closest('[hx-content]');
            if (element) {
                var triggerSpec = element.getAttribute("hx-content-trigger") || 'load';
                if (triggerSpec === 'change') {
                    var delay = parseFloat(element.getAttribute("hx-content-delay")) || 0;
                    if (delay > 0) {
                        setTimeout(function() {
                            this.processElement(element);
                        }.bind(this), delay * 1000);
                    } else {
                        this.processElement(element);
                    }
                }
            }
        },

        processElement: function (elt) {
            var contentAttr = elt.getAttribute('hx-content');
            if (contentAttr) {
                var instructions = contentAttr.split(',');
                instructions.forEach(function (instruction) {
                    var colonIndex = instruction.indexOf(':');
                    if (colonIndex === -1) return;

                    var selector = instruction.substring(0, colonIndex).trim();
                    var content = instruction.substring(colonIndex + 1).trim();
                    var target = document.querySelector(selector);

                    if (target) {
                        var swapType = elt.getAttribute("hx-content-swap") || 'innerHTML';
                        
                        if (swapType === 'innerHTML') {
                            target.innerHTML = content;
                        } else if (swapType === 'outerHTML') {
                            target.outerHTML = content;
                        } else if (swapType === 'append') {
                            target.insertAdjacentHTML('beforeend', content);
                        } else if (swapType === 'prepend') {
                            target.insertAdjacentHTML('afterbegin', content);
                        } else {
                            // Default to innerHTML if invalid swap type
                            target.innerHTML = content;
                        }
                    }
                });
            }
        },

        onEvent: function (name, evt) {
            if (name === 'htmx:afterSettle') {
                this.processLoadContentElements();
            }
        }
    });
})();