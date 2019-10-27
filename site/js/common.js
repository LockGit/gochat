let socketUrl = "ws://127.0.0.1:7000/ws";
let apiUrl = "http://127.0.0.1:7070";

function getLocalStorage(name) {
    return localStorage.getItem(name);
}

function setLocalStorage(name, value) {
    localStorage.setItem(name, value)
}