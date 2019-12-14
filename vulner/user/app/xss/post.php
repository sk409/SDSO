<?php
header("X-XSS-Protection: 0");
echo $_POST["name"];
echo htmlspecialchars($_POST["password"], ENT_QUOTES);
