export const FileUtility = {
  downloadFileFromBuffer(buffer, fileName) {
    const blob = new Blob([buffer], { type: "application/octet-stream" });
    const url = URL.createObjectURL(blob);

    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();

    URL.revokeObjectURL(url);
    document.body.removeChild(link);
  },

  handleFileInput(event, fileName, fileData) {
    event.preventDefault();

    const file = event.target.files[0];

    if (file === undefined) {
      fileName.value = undefined;
      fileData.value = undefined;
      return;
    }

    console.log(`selected file: ${file.name}`);
    fileName.value = file.name;
    fileData.value = file;
  },

  async fileToHex(file) {
    const file_bytes = new Uint8Array(await file.peek().arrayBuffer());

    const hex = Array.from(file_bytes)
      .map((b) => b.toString(16).padStart(2, "0"))
      .join("");

    return hex;
  },
};

export default FileUtility;
