function startTime() {
    const today = new Date();
    document.getElementById('datetime').innerHTML =  today.toLocaleString();
    setTimeout(startTime, 1000);
}
