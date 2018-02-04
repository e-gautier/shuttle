import urlInput from './urlInput';
import validationError from './validationError';
import $ from "jquery";
import {RCAPTCHA} from "./globals";

export default class shortenButton {
    constructor() {
        this.element    = $("#shortenButton");
        this.faLink     = $("#fa-link");
        this.faCog      = $("#fa-cog");
    }
    submitForm() {
        const u = new urlInput();
        const response = RCAPTCHA ? grecaptcha.getResponse() : null;

        $.ajax({
            url: '/',
            method: 'POST',
            data: {
                url: u.element.val(),
                response: response
            },
            success: (data) => {
                const v = new validationError();
                u.element.removeClass("is-invalid");
                v.element.empty();

                if (!data.error) {
                    this.showLinkIcon();
                    return u.updateInput(data.url);
                }

                u.element.addClass("is-invalid");
                $.each(data.error, function (i, e) {
                    v.element.append(v.render(e.Message));
                });

                this.showLinkIcon();
            },
            error: (xhr, statusText) => {
                const v = new validationError();
                v.element.append(v.render(statusText));
                this.showLinkIcon();
            }
        }).done(() => {
            RCAPTCHA ? grecaptcha.reset() : null;
        })
    }
    showLinkIcon() {
        this.faCog.css('display', 'none');
        this.faLink.css('display', 'inline');
    }
    showCogIcon() {
        this.faLink.css('display', 'none');
        this.faCog.css('display', 'inline');
    }
};
