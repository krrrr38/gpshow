<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8" />
    <title>{{ .Title }}</title>
    <link rel="stylesheet" type="text/css" href="{{ .Asset }}css/show.css" />
    <link rel="stylesheet" type="text/css" href="{{ .Prettyfy }}prettify.css" />
    <link rel="stylesheet" type="text/css" href="css/custom.css" />
    <script type="text/javascript" src="{{ .JQuery }}"></script>
    <script type="text/javascript" src="{{ .Asset }}js/show.js"></script>
    <script type="text/javascript" src="{{ .Prettyfy }}prettify.js"></script>
    {{ range $index, $lang := .Languages }}
      <script type="text/javascript" src="{{ $.Prettyfy }}lang-{{$lang}}.js"></script>
    {{ end }}
    <script type="text/javascript" src="js/custom.js"></script>
    <script type="text/javascript"><!--
      window.onload=function() { prettyPrint(); };
    --></script>
  </head>
  <body>
    <div id="slides">
      <div id="reel">
        {{ .Slides }}
      </div>
    </div>
  </body>
</html>
