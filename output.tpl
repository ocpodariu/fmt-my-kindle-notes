<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <title>Kindle Notes - {{ .Title }}</title>
</head>

<body>
    {{ range .Sections }}
        <div><b>{{ .Title }}</b></div>
        {{ range .Highlights }}
            <div><u>
                Highlight ({{ .Color }})
                {{ if .Page }} - Page {{ .Page }} {{ end }}
                {{ if .Location }} - Location {{ .Location }} {{ end }}
            </u></div>
            <div>&emsp;{{ .Text }}</div>
            {{ if .Note }}
                <div>&emsp;<i>{{ .Note }}</i></div>
            {{ end }}
        {{ end }}
        <br>
    {{ end }}
</body>

</html>
