let socketUrl = "ws://207.246.91.166:7000/ws";
let apiUrl = "http://207.246.91.166:7070";

function getLocalStorage(name) {
    let value = localStorage.getItem(name);
    if (value == null) {
        return "";
    }
    return value;
}

function setLocalStorage(name, value) {
    localStorage.setItem(name, value)
}
