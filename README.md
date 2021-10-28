# Photobucket Exodus
In the 2000s it was common to incorporate pictures in websites, forums and blogs by putting them in online sharing services such as Photobucket, ImageShack or TinyPic. Sadly most of these sites failed to find a profitable business model therefore they heaviliy reduced their free tier of ceased their operations altogether. This damaged lots of historical websites which suddenly most of their photos deleted or watermarked.

While deleted pictures from closed sites such as TinyPic are basically impossible to recover without a backup watermarked pictures from Photobucket can be downloaded on bulk with a properly configured script.

This little command line tool scan an input file (which can be a simple list, an HTML page, a SQL dump from a CMS...) locates all the links to Photobucket pictures, downloads them setting the right HTTP headers to trick the protection which Photobucket put in place to prevent this kind of massive downloads, reverts them to their original formats (Photobucket stores everything in webp) and in the end saves them on a folder with their original naming preserved making easy to upload them on another location and eventually change the links to the pictures on bulk (in fact in most cases a simple search & replace on an editor which supports regex such as Vim or Visual Studio Code is enough).

## WARNING!

While the tool is already usable it still lacks a very important feature which will be added in the next weeks: it cannot manage duplicates therefore currently if there are there are two pics with the same name in two different album the second picture will overwrite the first, please keep it in mind if you want to use this tool before the release of the version 1.0

