
// サーバーと通信を行う関数

// 現在再生中の曲を問い合わせる（他のデータもjsonで欲しいな...）
window.addEventListener("DOMContentLoaded", function () {
	
});

// 曲を変えるための関数（引数は仮のもの）
function control_Music(n){
	if(n == 1); // 次の曲
	if(n == -1); // リスタート
	if(n == 0);	// 停止
}

// 音量を切り替える
function control_Volume(n){
	if(n == -1); // 音量を下げる
	if(n == 1); // 音量をあげる
}

// プレイリスト一覧を描画する関数
function display_Playlist(s){
	
}

/*
	メモ
		曲の停止、再生とミュートはトグル？
/* */