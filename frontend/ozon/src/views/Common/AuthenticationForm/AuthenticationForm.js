import authenticationFormTemplate from './AuthenticationForm.hbs';
import buttonStyles from '../Button/Button.scss';
import textStyles from '../TextArea/TextArea.scss';
import linkStyles from '../Link/Link.scss';


/**
 * @class AuthenticationForm
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class AuthenticationForm {
    /**
     * A default constructor for AuthenticationForm
     */
    constructor() {
        this.objectType = 'authenticationForm';
    };

    /**
     *
     * @param {string} formType Type of templated Authentication form, which will generate
     * @param {Object} formData Array of objects which will rendered inside template
     * @return {string} HTML string, which will generated after templating
     */
    getHtmlString = ({formType, formData, styles}) => {
        return authenticationFormTemplate({
            formType: formType,
            formData: formData,
            styles: styles,
            buttonStyles: buttonStyles,
            textStyles: textStyles,
            linkStyles: linkStyles,
        });
    };
}

export default AuthenticationForm;
