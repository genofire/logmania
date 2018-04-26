import * as gui from './gui';
import config from './config';

/**
 * Self binding with router
 */
/* eslint-disable no-unused-vars */
import home from './view/home';
import login from './view/log';
/* eslint-enable no-unused-vars */


document.title = config.title;
window.onload = () => gui.render();
