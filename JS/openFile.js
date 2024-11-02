function openJsonFile() {
    const input = document.createElement('input');
    return new Promise((s, f) => {
        input.type = 'file';
        input.style.display = 'none';
        document.body.append(input);
        input.addEventListener('change', function(e) {
            const file = e.target.files[0];
            if (file) {
                const reader = new FileReader();

                reader.onload = function(e) {
                    try {
                        const content = e.target.result;
                        const data = JSON.parse(content);
                        // displayData(data);
                        s({
                            data,
                            name: file.name,
                        });
                    } catch (error) {
                        f("Error parsing JSON file:" + error.message);
                        // alert("无法解析JSON文件，请检查文件格式是否正确！");
                    }
                };

                reader.readAsText(file);
            } else {
                f("没有选择文件！");
            }
        });
        input.click();
    }).then(r => {
        input.remove();
        return r;
    }).catch(e => {
        input.remove();
        return e;
    });
}