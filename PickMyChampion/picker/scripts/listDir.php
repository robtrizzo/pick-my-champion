<?php
function getFiles($dirName){
    $files=array();
    if($dir=opendir($dirName)){
        while($file=readdir($dir)){
            if($file!='.' && $file!='..' && $file!=basename(__FILE__)){
                $files[]=$file;
            }
        }
        closedir($dir);
    }
    natsort($files); //sort
    return json_encode($files);
}

if (isset($_POST['dir_path'])) {
    echo getFiles($_POST['dir_path']);
}
console.log("calling getFiles")
?>
