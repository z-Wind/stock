# stock
TD API &amp; AlphaVantage API &amp; TWSE API

若要使用 gotd，請先放置 https 相關檔案，指令如下<br/>
<pre>
<code>
		openssl genrsa -out key.pem 2048
		openssl req -new -x509 -key key.pem -out cert.pem -days 3650
</code>
</pre>
cert.pem<br/>
key.pem<br/>
產生後，放到 gotd 資料夾下<br/>

初次啟動時，會出現 TD 認證頁面是正常的<br/>
因需獲取需可，保證資料皆存在你的電腦上<br/>

若有不需要的 api，可在 main.go 的 init 中修改，也可在此加入需要的 api<br/>

連線網址，會出現極度簡單的用法說明<br/>
<code>
	http://localhost:6060/
</code>
