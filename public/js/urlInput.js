import shortenButton from "./shortenButton";
import {RCAPTCHA} from "./globals";

export default class urlInput {
    constructor() {
        this.element = $("#urlInput");
    }
    updateInput(url) {
        const u = this.element;
        u.val(
            window.location.protocol+'//'+
            window.location.host+'/'+
            url
        );
        u.select();
        //document.execCommand("copy");
    }
    watch() {
        this.element.bind("paste", () => {
            setTimeout(() => {
                if (RCAPTCHA) {
                    grecaptcha.execute()
                } else {
                    const s = new shortenButton();
                    s.showCogIcon();
                    s.submitForm();
                }
            }, 100);
        });
    }
};
