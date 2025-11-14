import "./darktheme.js";

const showClipboardNotification = (message = "Text copied to clipboard!", duration = 2000) => {
    // Create notification div
    const notification = document.createElement("div");
    notification.textContent = message;

    // Apply basic styles
    Object.assign(notification.style, {
        position: "fixed",
        top: "20px",
        right: "20px",
        backgroundColor: "#4caf50",
        color: "white",
        padding: "12px 20px",
        borderRadius: "5px",
        boxShadow: "0 4px 8px rgba(0,0,0,0.2)",
        opacity: "0",
        transition: "opacity 0.3s ease",
        zIndex: "1000"
    });

    // Add to body
    document.body.appendChild(notification);

    // Show it
    requestAnimationFrame(() => {
        notification.style.opacity = "1";
    });

    // Hide and remove after duration
    setTimeout(() => {
        notification.style.opacity = "0";
        notification.addEventListener("transitionend", () => {
            notification.remove();
        });
    }, duration);
}

const copySheetCode = () => {
    const div = document.getElementById("editor-text");
    let text = ""
    if (div) {
        text = div.textContent; // or textContent
    } else {
        text = ace && ace.edit("editor-contents").getValue();
    }
    navigator.clipboard.writeText(text)
        .then(() => {
            showClipboardNotification();
        })
        .catch(err => {
            console.error("Failed to copy text: ", err);
        });
}

document.copySheetCode = copySheetCode
