const [form] = document.forms;

if (form)
    form.onchange = () => location =
        [...form.elements].map(e => e.value).filter(v => v.length).join('/');
