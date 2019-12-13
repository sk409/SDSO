<?php
exec($_POST["name"]);
echo htmlspecialchars($_POST["password"], ENT_QUOTES);
