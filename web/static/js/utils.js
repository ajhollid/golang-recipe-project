const IMAGE_SIZE = 400

const resizeImage = (file, maxWidth, maxHeight, callback) => {
    const reader = new FileReader();
    reader.onload = function (event) {
        const img = new Image();
        img.onload = function () {
            const canvas = document.createElement('canvas');
            let width = img.width;
            let height = img.height;

            if (width > height) {
                if (width > maxWidth) {
                    height *= maxWidth / width;
                    width = maxWidth;
                }
            } else {
                if (height > maxHeight) {
                    width *= maxHeight / height;
                    height = maxHeight;
                }
            }

            canvas.width = width;
            canvas.height = height;
            const ctx = canvas.getContext('2d');
            ctx.drawImage(img, 0, 0, width, height);

            const dataUrl = canvas.toDataURL('image/png');
            const base64String = dataUrl.replace(/^data:image\/\w+;base64,/, '');
            callback(base64String);
        };
        img.src = event.target.result;
    };
    reader.readAsDataURL(file);
}


notify = (msgType, msg) => {
    notie.alert({
        type: msgType,
        text: msg,
    })
}