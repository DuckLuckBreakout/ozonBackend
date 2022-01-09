import BaseView from '../BaseView.js';
import Input from '../Common/Input/Input.js';
import profileTemplate from './ProfileView.hbs';
import profileStyles from './ProfileView.scss';
import buttonStyles from '../Common/Button/Button.scss';
import textStyles from './../Common/TextArea/TextArea.scss';
import imgStyles from './../Common/Img/Img.scss';
import Events from '../../utils/bus/events';
import Img from '../Common/Img/Img';

/**
 * @class  ProfileView
 * @extends BaseView
 * @classdesc Class for showing profile of a user
 */
class ProfileView extends BaseView {
    /**
     * @param {Object} URLParams
     * @description redefinition of show method
     */
    show = (URLParams = {}) => {
        this.presenter.tryAuth();
    }

    /**
     *
     * @return {void} rendered page
     */
    render = () => {
        this.goUp();

        this.parent.innerHTML = '';
        if (this.cache !== '') {
            this.parent.appendChild(this.cache);
            this.bus.emit(Events.ProfileAllGet);
            return;
        }

        const htmlTemplate = profileTemplate({
            inputFields: [
                new Input({type: 'text', name: 'firstName', placeholder: 'Имя'}),
                new Input({type: 'text', name: 'lastName', placeholder: 'Фамилия'}),
            ],
            inputEmail: new Input({type: 'email', name: 'email',
                placeholder: 'Адрес электронной почты', isDisabled: true}),
            avatarUpload: new Input({type: 'file', name: 'avatar', placeholder: 'Upload new Avatar'}),
            profileStyles: profileStyles,
            buttonStyles: buttonStyles,
            imgStyles: imgStyles,
            textStyles: textStyles,
            smileSvg: new Img({src: '/svg/header/smile.svg'}),
        });

        this.cache = new DOMParser().parseFromString(htmlTemplate, 'text/html').getElementById('profile-page');
        this.parent.appendChild(this.cache);
        const form = document.getElementById('flname-form');
        form.addEventListener('submit', (evt) => {
            evt.preventDefault();
            this.bus.emit(Events.ProfileFLNameChange);
        });

        const avatarInput = document.getElementById('avatar-input');
        avatarInput.addEventListener('change', (evt) => {
            evt.preventDefault();
            this.bus.emit(Events.ProfileAvatarChange);
        });

        const logout = document.getElementById('logout');
        logout.addEventListener('click', (evt) => {
            evt.preventDefault();
            this.bus.emit(Events.ProfileLogoutPrepare);
        });

        this.bus.emit(Events.ProfileAllGet);
    }


    /**
     *
     * @param {string} firstName
     * @param {string} lastName
     */
    changeFirstLastName = (firstName, lastName) => {
        document.getElementsByName('firstName')[0].value = firstName;
        document.getElementsByName('lastName')[0].value = lastName;
        document.getElementById('profile-firstName').textContent = firstName;
        document.getElementById('profile-lastName').textContent = lastName;
    }

    /**
     *
     * @param {string} avatarURL
     */
    changeAvatar = (avatarURL) => {
        document.getElementById('avatar-img').src = avatarURL;
    }

    /**
     *
     * @param {string} email
     */
    changeEmail = (email) => {
        document.getElementsByName('email')[0].value = email;
    }

    /**
     * @description Using for render data after AJAX methods.
     */
    renderData = () => {
        this.changeFirstLastName(this.presenter.getFirstName(), this.presenter.getLastName());
        this.changeEmail(this.presenter.getEmail());
        this.changeAvatar(this.presenter.getAvatar());
    }

    invalidAvatar = () => {
        document.getElementById('invalid-avatar-span').innerHTML = 'Некорретный файл';
    }
}

export default ProfileView;
