const has = (input) => {
    return input.length > 0
}

const validAmount = (amount) => {
    const pattern = new RegExp("^\\d+\\/\\d+$|^\\d+(\\.\\d+)?$");
    return pattern.test(amount)
}

const minLength = (input, minLength) => {
    return input.length >= minLength
}