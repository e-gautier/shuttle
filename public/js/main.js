'use strict';

import {RCAPTCHA} from "./globals";

import 'babel-polyfill';
import $ from 'jquery';

window.jQuery = $;
window.$ = $;
global.jQuery = $;

import 'popper.js';
import 'bootstrap/dist/css/bootstrap.min.css';
import '../css/paper-bootstrap.min.css';
import '../css/main.css';

import fontawesome from '@fortawesome/fontawesome';
import faSpaceShuttle from '@fortawesome/fontawesome-free-solid/faSpaceShuttle';
import faLink from '@fortawesome/fontawesome-free-solid/faLink';
import faCog from '@fortawesome/fontawesome-free-solid/faCog';
fontawesome.library.add(faSpaceShuttle);
fontawesome.library.add(faLink);
fontawesome.library.add(faCog);

import 'particles.js';
import particles from '../assets/particles';
import urlInput from './urlInput';
import shortenButton from './shortenButton';

window.onSubmitWithCaptcha = function () {
    const s = new shortenButton();
    s.showCogIcon();
    s.submitForm();
};

function main() {
    particlesJS('particles-js', particles);
    $('form').submit((event) => {
        event.preventDefault();
        if (RCAPTCHA) {
            grecaptcha.execute()
        } else {
            const s = new shortenButton();
            s.showCogIcon();
            s.submitForm();
        }
    });
    const u = new urlInput();
    u.watch();
}

$(document).ready(() => {
    main();
});
