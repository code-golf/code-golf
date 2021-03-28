const [form] = document.forms;

if (form)
    form.onchange = () => location =
        [...new FormData(form).values()].filter(v => v.length).join('/');
