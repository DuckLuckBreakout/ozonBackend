/**
 * @param {Object} arg1 first argument to check for equality
 * @param {Object} arg2 second argument to check for equality
 * @param {Object} options Option object which method result will return
 * @return {Object} Object for hbs-loader which checks for true or false
 */
const ifEq = function(arg1, arg2, options) {
    return (arg1 === arg2) ? options.fn(this) : options.inverse(this);
};

export default ifEq;
