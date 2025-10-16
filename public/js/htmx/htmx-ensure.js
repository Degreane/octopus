// // (function () {
// //     htmx.defineExtension('htmx-ensure', {
// //         init: function (api) {
// //             var self = this;
// //             document.addEventListener('DOMContentLoaded', function () {
// //                 self.processLoadEnsureElements();
// //
// //                 // Listen for DOM changes to handle dynamically added elements
// //                 var observer = new MutationObserver(self.processLoadEnsureElements.bind(self));
// //                 observer.observe(document.body, {childList: true, subtree: true});
// //             });
// //
// //             // Add click event listener to document
// //             document.addEventListener('click', this.handleClick.bind(this));
// //         },
// //
// //         processLoadEnsureElements: function () {
// //             document.querySelectorAll('[hx-ensure-load]').forEach(this.processElement.bind(this, 'hx-ensure-load'));
// //         },
// //
// //         handleClick: function (event) {
// //             var element = event.target.closest('[hx-ensure]');
// //             if (element) {
// //                 this.processElement('hx-ensure', element);
// //             }
// //         },
// //
// //         processElement: function (attrName, elt) {
// //             var ensureAttr = elt.getAttribute(attrName);
// //             if (ensureAttr) {
// //                 var instructions = ensureAttr.split(',');
// //                 instructions.forEach(function (instruction) {
// //                     var [selector, actions] = instruction.trim().split(':');
// //                     var target = document.querySelector(selector);
// //
// //                     if (target) {
// //                         actions.split(';').forEach(function (action) {
// //                             if (action.includes('s')) {
// //                                 // Timed action
// //                                 var [className, time] = action.split('s');
// //                                 var remove = className.startsWith('!');
// //                                 className = remove ? className.substring(1) : className;
// //
// //                                 setTimeout(function () {
// //                                     if (remove) {
// //                                         target.classList.remove(className);
// //                                     } else {
// //                                         target.classList.add(className);
// //                                     }
// //                                 }, parseInt(time) * 1000);
// //                             } else {
// //                                 // Immediate action
// //                                 var remove = action.startsWith('!');
// //                                 var className = remove ? action.substring(1) : action;
// //
// //                                 if (remove) {
// //                                     target.classList.remove(className);
// //                                 } else {
// //                                     target.classList.add(className);
// //                                 }
// //                             }
// //                         });
// //                     }
// //                 });
// //             }
// //         },
// //
// //         onEvent: function (name, evt) {
// //             if (name === 'htmx:afterSettle') {
// //                 this.processLoadEnsureElements();
// //             }
// //         }
// //     });
// //
// // })();
// //
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
//                             var act = action.trim();
//                             // Support formats:
//                             //  - "className"          => add immediately
//                             //  - "!className"         => remove immediately
//                             //  - "className 3s"       => add after 3s
//                             //  - "!className 1.5s"    => remove after 1.5s
//                             var timedMatch = act.match(/^(!?)(\S+)\s+(\d+(?:\.\d+)?)s$/);
//
//                             if (timedMatch) {
//                                 var removeTimed = timedMatch[1] === '!';
//                                 var classNameTimed = timedMatch[2];
//                                 var seconds = parseFloat(timedMatch[3]);
//
//                                 setTimeout(function () {
//                                     if (removeTimed) {
//                                         target.classList.remove(classNameTimed);
//                                     } else {
//                                         target.classList.add(classNameTimed);
//                                     }
//                                 }, seconds * 1000);
//                             } else {
//                                 var remove = act.startsWith('!');
//                                 var className = remove ? act.substring(1) : act;
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
(function () {
    htmx.defineExtension('htmx-ensure', {
        init: function (api) {
            var self = this;
            document.addEventListener('DOMContentLoaded', function () {
                self.processLoadEnsureElements();
                self.bindEnsureTriggers();

                // Listen for DOM changes to handle dynamically added elements
                var observer = new MutationObserver(function () {
                    self.processLoadEnsureElements();
                    self.bindEnsureTriggers();
                });
                observer.observe(document.body, {childList: true, subtree: true});
            });

            // Fallback: delegated click handler (in case hx-trigger is omitted)
            document.addEventListener('click', this.handleClick.bind(this));
        },

        // Ensure actions defined via hx-ensure-load run on load/settle
        processLoadEnsureElements: function () {
            document.querySelectorAll('[hx-ensure-load]').forEach(this.processElement.bind(this, 'hx-ensure-load'));
        },

        // Backward-compatible click behavior when hx-trigger not specified
        handleClick: function (event) {
            var element = event.target.closest('[hx-ensure]');
            if (element && !element.hasAttribute('hx-trigger')) {
                this.processElement('hx-ensure', element);
            }
        },

        // NEW: Bind event listeners according to hx-trigger on [hx-ensure] elements
        bindEnsureTriggers: function () {
            var self = this;
            var ensureNodes = document.querySelectorAll('[hx-ensure]');
            ensureNodes.forEach(function (elt) {
                // Prevent duplicate binding on the same node
                if (elt.dataset.ensureBound === '1') return;

                var triggers = (elt.getAttribute('hx-trigger') || 'click')
                    .split(',')
                    .map(function (t) { return t.trim().split(/\s+/)[0]; }) // strip simple modifiers if present
                    .filter(Boolean);

                // If 'load' specified, schedule immediate processing for this element
                if (triggers.includes('load')) {
                    self.scheduleLoadEnsure(elt);
                }

                // Bind non-load events (e.g., click, change, input, mouseenter, etc.)
                triggers
                    .filter(function (t) { return t.toLowerCase() !== 'load'; })
                    .forEach(function (evName) {
                        elt.addEventListener(evName, function () {
                            self.processElement('hx-ensure', elt);
                        }, { passive: true });
                    });

                // Mark as bound once we attach listeners or schedule load
                elt.dataset.ensureBound = '1';
            });
        },

        // NEW: Run hx-ensure for elements whose hx-trigger contains 'load'
        scheduleLoadEnsure: function (elt) {
            // Avoid multiple fires for the same node across DOM updates
            if (elt.dataset.ensureLoadFired === '1') return;

            var self = this;
            // Run after insertion/rendering; requestAnimationFrame ensures layout is ready
            var run = function () {
                // Double-check still in DOM
                if (document.contains(elt)) {
                    self.processElement('hx-ensure', elt);
                    elt.dataset.ensureLoadFired = '1';
                }
            };

            if (document.readyState === 'complete' || document.readyState === 'interactive') {
                requestAnimationFrame(run);
            } else {
                window.addEventListener('load', function onWinLoad() {
                    window.removeEventListener('load', onWinLoad);
                    requestAnimationFrame(run);
                });
            }
        },

        processElement: function (attrName, elt) {
            var ensureAttr = elt.getAttribute(attrName);
            if (ensureAttr) {
                var instructions = ensureAttr.split(',');
                instructions.forEach(function (instruction) {
                    var parts = instruction.trim().split(':');
                    var selector = parts[0];
                    var actions = parts.slice(1).join(':'); // allow ':' inside actions safely
                    var target = document.querySelector(selector);

                    if (target && actions) {
                        actions.split(';').forEach(function (action) {
                            var act = action.trim();
                            // Support formats:
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
                            } else if (act.length) {
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
            // Re-process after htmx swaps and newly settled content
            if (name === 'htmx:afterSettle' || name === 'htmx:load') {
                this.processLoadEnsureElements();
                this.bindEnsureTriggers();

                // Fire load triggers for newly inserted elements that have not fired yet
                document.querySelectorAll('[hx-ensure][hx-trigger*="load"]').forEach(this.scheduleLoadEnsure.bind(this));
            }
        }
    });

})();