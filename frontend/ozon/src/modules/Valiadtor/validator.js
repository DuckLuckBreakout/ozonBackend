import {fields} from '../../utils/validationFields/validationFields';

export const isValidInputs = (inputs, specificTypesToCheck = []) => {
    const userInfo = {
        password: '',
    };

    const result = {
        failedFields: [],
    };

    for (const input of inputs) {
        if (specificTypesToCheck.length > 0 &&
           !specificTypesToCheck.includes(input.type)
        ) {
            continue;
        }

        if (input.disabled) {
            continue;
        }

        switch (input.type) {
        case 'text':
            if (input.name.includes('firstName') ||
                input.name.includes('lastName')
            ) {
                nameValidation(input, result);
            }
            break;
        case 'email':
            emailValidation(input, result);
            break;
        case 'file':
            fileValidation(input, result);
            break;
        case 'password':
            switch (input.id) {
            case 'password':
                passwordValidation(input, result, userInfo);
                break;
            case 'repeatPassword':
                passwordRepeatValidation(input, result, userInfo.password);
                break;
            default:
                break;
            }
            break;

        default:
            console.error('unsupported type of input!');
            break;
        }
    }

    return result;
};

const nameValidation = (input, result) => {
    const isValidName = /^[a-zA-Zа-яА-Я]+$/.test(input.value);
    if (!isValidName) {
        switch (input.name) {
        case 'firstName':
            result.failedFields.push(fields.incorrectFirstName);
            break;
        case 'lastName':
            result.failedFields.push(fields.incorrectLastName);
            break;
        }
    }
};

const emailValidation = (input, result) => {
    const isValidEmail = /^([a-z0-9_-]+\.)*[a-z0-9_-]+@[a-z0-9_-]+(\.[a-z0-9_-]+)*\.[a-z]{2,6}$/.test(input.value);
    if (!isValidEmail) {
        result.failedFields.push(fields.incorrectEmail);
    }
};

const passwordValidation = (input, result, userInfo) => {
    if (input.value.length < 3) {
        result.failedFields.push(fields.passwordTooShort);
        return;
    }
    if (input.value.length > 20) {
        result.failedFields.push(fields.passwordTooLong);
        return;
    }
    userInfo.password = input.value;
};

const passwordRepeatValidation = (input, result, password) => {
    if (result.failedFields.includes(fields.passwordTooLong) ||
        result.failedFields.includes(fields.passwordTooShort)) {
        result.failedFields.push(fields.repeatPasswordDrop);
        return;
    }

    const value = input.value;
    if (value !== password) {
        result.failedFields.push(fields.repeatPasswordUnmatched);
        return;
    }
    if (input.value.length < 3) {
        result.failedFields.push(fields.repeatPasswordTooShort);
        return;
    }
    if (input.value.length > 20) {
        result.failedFields.push(fields.repeatPasswordTooLong);
    }
};

const MAX_FILE_SIZE = 10e6;
const fileValidation = (input, result) => {
    const file = input.files[0];
    const isValidFile = file !== undefined &&
                         file.size < MAX_FILE_SIZE && /.*\.(jpeg|png|jpg)$/i.test(file.name);
    if (!isValidFile) {
        result.failedFields.push(fields.incorrectFile);
    }
};
