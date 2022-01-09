import BaseView from '../BaseView.js';
import Input from '../Common/Input/Input.js';
import Button from '../Common/Button/Button.js';
import Link from '../Common/Link/Link.js';
import Popup from '../Common/Popup/Popup.js';
import Blind from '../Common/Blind/Blind.js';
import AuthenticationForm from '../Common/AuthenticationForm/AuthenticationForm.js';
import AuthenticationFormStyles from '../Common/AuthenticationForm/AuthenticationForm.scss';
import Router from '../../utils/router/Router';
import Events from '../../utils/bus/events';
import decorator from '../decorators.scss';
import buttonStyles from '../Common/Button/Button.scss';
import popupStyles from '../Common/Popup/Popup.scss';
import linkStyles from '../Common/Link/Link.scss';

/**
 * @class LoginView
 * @extends BaseView
 * @classdesc Class for Login page
 */
class LoginView extends BaseView {
    /**
     *
     * @return {void} html form
     */
    render = () => {
        const body = document.getElementsByTagName('body')[0];
        body.classList.add(decorator.noScroll);

        if (this.cache !== '') {
            const inputs = this.cache.getElementsByTagName('input');
            inputs[0].placeholder = 'Электронная почта';
            inputs[1].placeholder = 'Пароль';
            for (const input of inputs) {
                input.style['border-color'] = '';
            }

            this.parent.appendChild(this.cache);
            return;
        }

        const template = new Popup().getHtmlString({
            popupBody:
                new AuthenticationForm().getHtmlString({
                    formType: 'login',
                    formData: {
                        inputFields: [
                            new Input({
                                type: 'email',
                                name: 'email',
                                placeholder: 'Электронная почта',
                            }),
                            new Input({
                                type: 'password',
                                name: 'password',
                                placeholder: 'Пароль',
                            }),
                        ],
                        link: new Link({
                            href: '/signup',
                            name: 'Зарегистрироваться',
                        }),
                        submitButton: new Button({
                            text: 'Войти',
                            type: 'submit',
                        }),
                    },
                    styles: AuthenticationFormStyles,
                }),
            background: new Blind().getHtmlString(),
            popupType: popupStyles.login,
        });
        this.cache = new DOMParser().parseFromString(template, 'text/html').getElementById('popup');
        this.parent.appendChild(this.cache);

        document.getElementById('blind')
            .addEventListener('click', (evt) => {
                evt.preventDefault();
                body.classList.remove(decorator.noScroll);
                document.getElementById('popup').remove();
                Router.goBack();
            });

        const form = this.cache.getElementsByClassName(buttonStyles.auth)[0];
        form.addEventListener('click', (evt) => {
            body.classList.remove(decorator.noScroll);
            evt.preventDefault();
            this.bus.emit(Events.LoginSendData);
        });


        this.cache.getElementsByClassName(linkStyles.link)[0]
            .addEventListener('click', (evt) => {
                evt.preventDefault();
                body.classList.remove(decorator.noScroll);
                document.getElementById('popup').remove();
                Router.open('/signup', {replaceState: true});
            });
    }

    incorrectEmailOrPassword = () => {
        for (const input of this.cache.getElementsByTagName('input')) {
            BaseView.setInvalidInputPlaceholder(input, 'Неверные входные данные');
        }
    }
}

export default LoginView;
