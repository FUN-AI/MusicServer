
// サーバーと通信を行う関数

// 現在再生中の曲を問い合わせる（他のデータもjsonで欲しいな...）
window.addEventListener("DOMContentLoaded", function () {
	display_Playlist(["データ1","データ2"]);
});

// 曲を変えるための関数（引数は仮のもの）
function control_Music(n) {
	if (n == 1); // 次の曲
	if (n == -1); // リスタート
	if (n == 0);	// 停止
}

// 音量を切り替える
function control_Volume(n) {
	if (n == -1); // 音量を下げる
	if (n == 1); // 音量をあげる
}

// プレイリスト一覧を描画する関数
function display_Playlist(s) {
	// PlayList--Frameのなかにリストを入れる
	e = document.getElementById("PlayList--Frame");
	while (e.firstChild) e.removeChild(e.firstChild);
	if (s == null || s == [""]) {
		var li = document.createElement("li");
		li.textContent = "プレイリストが存在しません。";
		e.appendChild(li);
	} else {
		for(var i in s) {
			var li = document.createElement("li");
			var a = document.createElement("a");
			a.textContent = s[i];
			a.setAttribute("href", "javascript:void(0);");
			a.setAttribute("onclick", `load_Playlsit("${s[i]}")`);
			li.appendChild(a);
			e.appendChild(li);
		}
	}
}

// プレイリストをロードする
function load_Playlsit(s){

}

/*
	メモ
		曲の停止、再生とミュートはトグル？
/* */