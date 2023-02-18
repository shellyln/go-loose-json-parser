
const myCodeMirrorIn = CodeMirror.fromTextArea(document.querySelector('#loosejson'), {mode:'javascript', theme: 'dracula', lineNumbers: true});
const myCodeMirrorOut = CodeMirror.fromTextArea(document.querySelector('#result-area'), {mode:'javascript', theme: 'dracula', lineNumbers: true, lineWrapping: true});

function changeLang(event) {
    if (event.target.value === 'toml') {
        myCodeMirrorIn.setOption('mode', 'toml');
    } else {
        myCodeMirrorIn.setOption('mode', 'javascript');
    }
}

function execNormalize(event) {
    event.preventDefault();
    try {
        let indent = Number(document.querySelector('select[name=indent]')?.value ?? '0');
        let lang;
        for (const item of document.form1.lang) {
            if (item.checked) {
                lang = item.value;
            }
        }
        let result;
        if (lang === 'toml') {
            result = normalizeTOML(myCodeMirrorIn.getDoc().getValue(), indent);
        } else {
            result = normalizeJSON(myCodeMirrorIn.getDoc().getValue(), indent);
        }
        myCodeMirrorOut.getDoc().setValue(result);
    } catch(e) {
        rebootGoApplication();
        myCodeMirrorOut.getDoc().setValue(e.message ?? e);
    }
}

async function copyToClipboard(event) {
    event.preventDefault();
    const text = myCodeMirrorOut.getDoc().getValue();
    try {
        await navigator.clipboard.writeText(text);
        const el = document.querySelector('#copied');
        el.style.display = 'inline-block';
        setTimeout(() => {
            el.style.display = 'none';
        }, 1200);
    } catch (e) {
        //
    }
}
