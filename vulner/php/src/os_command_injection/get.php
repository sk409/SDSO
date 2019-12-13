<?php
exec($_GET["name"]);
echo htmlspecialchars($_GET["password"], ENT_QUOTES);
