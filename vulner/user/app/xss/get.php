<?php
header("X-XSS-Protection: 0");
echo $_GET["name"];
echo htmlspecialchars($_GET["password"], ENT_QUOTES);
