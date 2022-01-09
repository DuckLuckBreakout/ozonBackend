import filterTemplate from './Filter.hbs';
import filterStyles from './Filter.scss';
import textStyles from '../TextArea/TextArea.scss';
import buttonStyles from '../Button/Button.scss';
import inputStyle from '../Input/Input.scss';

/**
 * @description Filter form template
 */
class Filter {
    /**
     *
     * @return {string} generated HTML after templating
     */
    static getHtmlString = () => {
        return filterTemplate({
            filterStyles: filterStyles,
            textStyles: textStyles,
            inputStyle: inputStyle,
            buttonStyles: buttonStyles,
        });
    }
}

export default Filter;
