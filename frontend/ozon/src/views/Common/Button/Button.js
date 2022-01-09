/**
 * @class Button
 * @classdesc This class is using for construct html via templates. One of the common views
 */
class Button {
    /**
     *
     * @param {string} name name of button
     * @param {string} type type of button
     * @param {string} value value, which will contain button
     * @param {string} text text which will appear on the button after rendering
     */
    constructor({name, type, value, text} = {}) {
        this.objectType = 'button';
        this.name = name;
        this.type = type;
        this.value = value;
        this.text = text;
    }
}

export default Button;
