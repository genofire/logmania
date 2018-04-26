import * as domlib from '../domlib';
import * as gui from '../gui';
import * as socket from '../socket';
import View from '../view';
import {store} from '../store';


function levelToColor (lvl) {
	return lvl;
}

function addItem (el, msg) {
	const div = domlib.newAt(el, 'div', {
		'class': levelToColor(msg.Leve)
	});
	domlib.newAt(div, 'span', null, msg.Data.hostname);
	domlib.newAt(div, 'span', null, msg.Message);
}

class LogView extends View {

	// eslint-disable-next-line class-methods-use-this
	render () {
		if (!this.init) {
			this.init = true;
		}

		/*
		 * Domlib.newAt(this.el, 'h2', {'class': 'ui header'}, 'Log');
		 * for (const msg in store.channel.ffhb) {
		 * domlib.newAt(this.el, 'div', null, msg.Data.hostname, msg);
		 *}
		 */
	}


	constructor () {
		super();
		socket.addEvent('ws:', (msg) => {
			// Length('ws:') = 3
			// eslint-disable-next-line no-magic-numbers
			const channel = msg.subject.substr(3);
			if (!store.channel[channel]) {
				store.channel[channel] = [];
			}
			store.channel[channel].push(msg.body);
			addItem(this.el, msg.body);
			this.render();
		});
	}
}

const logView = new LogView();

gui.router.on('/log', () => {
	gui.setView(logView);
});
