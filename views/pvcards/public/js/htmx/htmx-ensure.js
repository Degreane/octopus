// (function () {
//     htmx.defineExtension('htmx-ensure', {
//         init: function (api) {
//             var self = this;
//             document.addEventListener('DOMContentLoaded', function () {
//                 self.processLoadEnsureElements();
//
//                 // Listen for DOM changes to handle dynamically added elements
//                 var observer = new MutationObserver(self.processLoadEnsureElements.bind(self));
//                 observer.observe(document.body, {childList: true, subtree: true});
//             });
//
//             // Add click event listener to document
//             document.addEventListener('click', this.handleClick.bind(this));
//         },
//
//         processLoadEnsureElements: function () {
//             document.querySelectorAll('[hx-ensure-load]').forEach(this.processElement.bind(this, 'hx-ensure-load'));
//         },
//
//         handleClick: function (event) {
//             var element = event.target.closest('[hx-ensure]');
//             if (element) {
//                 this.processElement('hx-ensure', element);
//             }
//         },
//
//         processElement: function (attrName, elt) {
//             var ensureAttr = elt.getAttribute(attrName);
//             if (ensureAttr) {
//                 var instructions = ensureAttr.split(',');
//                 instructions.forEach(function (instruction) {
//                     var [selector, actions] = instruction.trim().split(':');
//                     var target = document.querySelector(selector);
//
//                     if (target) {
//                         actions.split(';').forEach(function (action) {
//                             if (action.includes('s')) {
//                                 // Timed action
//                                 var [className, time] = action.split('s');
//                                 var remove = className.startsWith('!');
//                                 className = remove ? className.substring(1) : className;
//
//                                 setTimeout(function () {
//                                     if (remove) {
//                                         target.classList.remove(className);
//                                     } else {
//                                         target.classList.add(className);
//                                     }
//                                 }, parseInt(time) * 1000);
//                             } else {
//                                 // Immediate action
//                                 var remove = action.startsWith('!');
//                                 var className = remove ? action.substring(1) : action;
//
//                                 if (remove) {
//                                     target.classList.remove(className);
//                                 } else {
//                                     target.classList.add(className);
//                                 }
//                             }
//                         });
//                     }
//                 });
//             }
//         },
//
//         onEvent: function (name, evt) {
//             if (name === 'htmx:afterSettle') {
//                 this.processLoadEnsureElements();
//             }
//         }
//     });
//
// })();
//
(function () {
    htmx.defineExtension('htmx-ensure', {
        init: function (api) {
            var self = this;
            document.addEventListener('DOMContentLoaded', function () {
                self.processLoadEnsureElements();

                // Listen for DOM changes to handle dynamically added elements
                var observer = new MutationObserver(self.processLoadEnsureElements.bind(self));
                observer.observe(document.body, {childList: true, subtree: true});
            });

            // Add click event listener to document
            document.addEventListener('click', this.handleClick.bind(this));
        },

        processLoadEnsureElements: function () {
            document.querySelectorAll('[hx-ensure-load]').forEach(this.processElement.bind(this, 'hx-ensure-load'));
        },

        handleClick: function (event) {
            var element = event.target.closest('[hx-ensure]');
            if (element) {
                this.processElement('hx-ensure', element);
            }
        },

        processElement: function (attrName, elt) {
            var ensureAttr = elt.getAttribute(attrName);
            if (!ensureAttr) {
                return;
            }

            ensureAttr.split(',').forEach(function (instructionRaw) {
                var instruction = (instructionRaw || '').trim();
                if (!instruction) {
                    return;
                }

                var colonIdx = instruction.indexOf(':');
                if (colonIdx === -1) {
                    return;
                }

                var selector = instruction.slice(0, colonIdx).trim();
                var actions = instruction.slice(colonIdx + 1).trim();
                if (!selector || !actions) {
                    return;
                }

                var targets;
                try {
                    targets = document.querySelectorAll(selector);
                } catch (e) {
                    return;
                }
                if (!targets || targets.length === 0) {
                    return;
                }

                actions.split(';').forEach(function (actionRaw) {
                    var act = (actionRaw || '').trim();
                    if (!act) {
                        return;
                    }

                    var timedMatch = act.match(/^(!?)(\S+)\s+(\d+(?:\.\d+)?)s$/);

                    if (timedMatch) {
                        var removeTimed = timedMatch[1] === '!';
                        var classNameTimed = timedMatch[2];
                        var seconds = parseFloat(timedMatch[3]);

                        targets.forEach(function (target) {
                            setTimeout(function () {
                                if (removeTimed) {
                                    target.classList.remove(classNameTimed);
                                } else {
                                    target.classList.add(classNameTimed);
                                }
                            }, seconds * 1000);
                        });
                    } else {
                        var remove = act.startsWith('!');
                        var className = remove ? act.substring(1) : act;

                        targets.forEach(function (target) {
                            if (remove) {
                                target.classList.remove(className);
                            } else {
                                target.classList.add(className);
                            }
                        });
                    }
                });
            });
        },

        onEvent: function (name, evt) {
            if (name === 'htmx:afterSettle') {
                this.processLoadEnsureElements();
            }
        }
    });

})();