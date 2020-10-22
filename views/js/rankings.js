const [form] = document.forms;

form.onchange = () => location = [...form.elements].map(e => e.value).join('/');
