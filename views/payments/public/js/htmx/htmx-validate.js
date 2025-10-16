(function () {
    htmx.on('htmx:afterProcess', function (event) {
        const element = event.detail.elt;
        const validateAttributes = Array.from(element.attributes).filter(attr => attr.name.startsWith('hx-validate'));

        validateAttributes.forEach(attr => {
            const regexString = attr.value;
            try {
                const regex = new RegExp(regexString);
                const value = element.value;
                const isValid = regex.test(value);

                const validTarget = element.getAttribute('hx-valid-target');
                const invalidTarget = element.getAttribute('hx-invalid-target');

                if (isValid) {
                    if (validTarget) {
                        document.querySelector(validTarget).style.display = 'block'; // or other action
                    }
                    if (invalidTarget) {
                        document.querySelector(invalidTarget).style.display = 'none'; // or other action
                    }
                } else {
                    if (validTarget) {
                        document.querySelector(validTarget).style.display = 'none'; // or other action
                    }
                    if (invalidTarget) {
                        document.querySelector(invalidTarget).style.display = 'block'; // or other action
                    }
                }
            } catch (error) {
                console.error("Invalid regular expression:", regexString, error);
            }
        });
    });

    // Optional: Add event listeners for dynamic validation (e.g., input change)
    htmx.on('htmx:afterSwap', function (event) {
        const elements = event.detail.elt.querySelectorAll('[hx-ext="htmx-validate"]');
        elements.forEach(el => {
            htmx.trigger(el, 'htmx:afterProcess', { elt: el }); //Trigger validation after swap
        });
    });

    //Handle trigger attributes for dynamic validation
    document.addEventListener('htmx:afterProcess', function (event) {
        const element = event.detail.elt;
        if (element.hasAttribute('hx-validate-trigger')) {
            const triggers = element.getAttribute('hx-validate-trigger').split(',').map(t => t.trim());
            triggers.forEach(trigger => {
                const [eventName, delay] = trigger.split(' ').map(t => t.trim());
                let timeoutId;
                element.addEventListener(eventName, function (e) {
                    clearTimeout(timeoutId);
                    if (delay) {
                        timeoutId = setTimeout(() => {
                            htmx.trigger(element, 'htmx:afterProcess', { elt: element });
                        }, parseFloat(delay) * 1000);
                    } else {
                        htmx.trigger(element, 'htmx:afterProcess', { elt: element });
                    }
                });
            });
        }
    });
})();


/*

(function () {
    htmx.on('htmx:afterProcess', function (event) {
        const element = event.detail.elt;
        const validateAttributes = Array.from(element.attributes).filter(attr => attr.name.startsWith('hx-validate'));

        validateAttributes.forEach(attr => {
            let isValid = false;
            const regexString = attr.value;
            const validationFunction = element.getAttribute('hx-validate-func');


            try {
                if (validationFunction) {
                    //Attempt to get the function by name if specified
                    const func = window[validationFunction];
                    if (typeof func === 'function') {
                        isValid = func(element.value);
                    } else {
                        console.error(`Validation function "${validationFunction}" not found.`);
                    }
                } else if (regexString) {
                    // Fallback to regex validation if no function is specified
                    const regex = new RegExp(regexString);
                    isValid = regex.test(element.value);
                }

                const validTarget = element.getAttribute('hx-valid-target');
                const invalidTarget = element.getAttribute('hx-invalid-target');

                if (isValid) {
                    // ... (rest of your existing code for handling valid input)
                } else {
                    // ... (rest of your existing code for handling invalid input)
                }
            } catch (error) {
                console.error("Validation error:", error);
            }
        });
    });

    // ... (rest of your existing code for htmx:afterSwap and dynamic validation)
})();

*/