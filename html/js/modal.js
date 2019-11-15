function load_Page(pass) {
	e = document.getElementById("modal");
	document.getElementById("modal_Title").textContent = "Now loading...";
	ef = document.getElementById("modal_Content");
	while (ef.firstChild) ef.removeChild(e.firstChild);
	e.style.display = "block";

	// ここからロードの処理
}