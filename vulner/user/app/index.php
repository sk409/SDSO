<?php

// class QueryItem
// {
//     private $name = "";
//     private $value = "";
// }

class Page
{
    public $title = "";
    public $getURL = "";
    public $postURL = "";
    public $inputs = "";

    public function __construct(string $title, string $url, array $queryItems)
    {
        $this->title = $title;
        $this->getURL = $url . "/get.php?";
        $this->postURL = $url . "/post.php";
        foreach ($queryItems as $key => $value) {
            $this->getURL .= ($key . "=" . $value . "&");
            $this->inputs .= "<input type=\"hidden\" name=\"$key\" value=\"$value\">";
        }
    }
}

$pages = [];
$pages[] = new Page("XSS", "xss", ["name" => "<script>alert(1);</script>", "password" => "pass"]);
$pages[] = new Page("OSコマンドインジェクション", "os_command_injection", ["name" => "sleep 1", "password" => "pass"]);
$pages[] = new Page("SQLインジェクション(MySQL-1)", "sql_injection/mysql/1", ["name" => "select sleep(1)", "password" => "pass"]);
$pages[] = new Page("SQLインジェクション(MySQL-2)", "sql_injection/mysql/2", ["name" => "'';select sleep(1)", "password" => "pass"]);
$pages[] = new Page("SQLインジェクション(MySQL-3)", "sql_injection/mysql/3", ["name" => "');select sleep(1);#", "password" => "pass"]);
?>

<!DOCTYPE html>
<html lang="ja">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Vulner</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <style>
        .container {
            display: grid;
            grid-template-columns: 30% 30% 30%;
            grid-column-gap: 3%;
            grid-row-gap: 3rem;
        }
    </style>
</head>

<body class="p-5">
    <div class="container">
        <?php
        foreach ($pages as $page) {
            echo <<<EOM
<div class="card">
    <h6 class="card-header">$page->title</h6>
    <div class="p-5 text-center">
        <a class="btn btn-primary" href="$page->getURL">GET</a>
        <form class="mt-3" method="post" action="$page->postURL">
            $page->inputs
            <input class="btn btn-primary" type="submit" value="POST">
        </form>
    </div>
</div>
EOM;
        }
        ?>
        <!-- <div class="card">
            <h5 class="card-header">XSS</h5>
            <div class="p-5 text-center">
                <a class="btn btn-primary" href="xss/get.php?name=<script>alert(1);</script>&password=def">GET</a>
                <form class="mt-3" method="post" action="xss/post.php">
                    <input type="hidden" name="name" value="<script>alert('name');</script>">
                    <input type="hidden" name="password" value="<script>alert('password');</script>">
                    <input class="btn btn-primary" type="submit" value="POST">
                </form>
            </div>
        </div>
        <div class="card">
            <h5 class="card-header">OSコマンドインジェクション</h5>
            <div class="p-5 text-center">
                <a class="btn btn-primary" href="os_command_injection/get.php?name=sleep 1&password=abc">GET</a>
                <form class="mt-3" method="post" action="os_command_injection/post.php">
                    <input type="hidden" name="name" value="sleep 1">
                    <input type="hidden" name="password" value="abc">
                    <input class="btn btn-primary" type="submit" value="POST">
                </form>
            </div>
        </div>
        <div class="card">
            <h5 class="card-header">SQLインジェクション(MySQL-1)</h5>
            <div class="p-5 text-center">
                <a class="btn btn-primary" href="sql_injection/mysql/get.php?name='';select sleep(5);&password=abc">GET</a>
                <form class="mt-3" method="post" action="sql_injection/mysql/post.php">
                    <input type="hidden" name="name" value="'';select sleep(5);">
                    <input type="hidden" name="password" value="abc">
                    <input class="btn btn-primary" type="submit" value="POST">
                </form>
            </div>
        </div> -->
    </div>
</body>

</html>