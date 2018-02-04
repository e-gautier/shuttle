export default class validationError {
    constructor() {
        this.element = $("#validation-container");
    }
    render(message) {
        return `<div style="color: red" id="validation">${message}</div>`;
    }
};
