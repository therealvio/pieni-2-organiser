# pieni-2-organiser

A dumb file renamer for [a dumber camera](https://kenkoglobal.com/product/toy_camera_pieni_ii/).

<!-- ![Product Image](./docs/product_image.png) -->
<img src="./docs/product_image.png" style="width:80%;display:block;margin-left:auto;margin-right:auto">

<div align="center"><em>I AM GOING TO SHIT</em></div>

## Usage

> [!NOTE]
> As of writing this, I am not publishing releases, so you will need to have `go` installed on your machine and compile/run it yourself.

This steaming pile of code looks over a directory recursively and creates copies of files in an output directory using a sequentual order, with sequencing grouped by the file extension.

Assuming you have all the media taken with this fuckass camera in one parent folder that you want to organise into one collection, run the below command:

```shell
go run main.go $directory_where_you_want_to_organise_photos
```

## Why?

This camera names its files in order of creation for that specific type. Without an onboard clock, or a way to even set one, this camera is as dumb as they come. Meaning, if you use multiple SD cards to store your media for a specific event, you will have multiple files with the same name. Which makes it hard if you want to put them together into one folder.

Let's be honest, needless folder structure is a dogshit experience 😡 and renaming files manually is equally a dogshit experience.

So, this thing exists to make renaming these files a little more humane. You don't have to use it, but if you do: you're welcome 😤.

To be clear: this was written for my own uses, but I thought someone else may have this camera or a similar use case, so I why not share the code?

## Lack of license?

Ah, so you've noticed. Yes, it is ambiguous whether or not you can use it. But do not worry: you as the individual can use it, clone it and modify it for your needs, I am more than okay with this. Just give me credit at least!

As for [Kenko](https://kenkoglobal.com/) and any other entity that makes more money than me[^1] will need to pony up or get in touch with me on using this or reselling it.

## Todo

No Github Issues for me, I am making this Github repo _work for me_ not _work for it_. Like I mentioned earlier, this repo is mostly for my own use-case, so I may or may not ever get around to do these items, but I at least thought about them.

- Use a CLI framework to provide enhanced inline help;
- Setup input validation for paths;
- Generate an output folder named after the input folder;
- Support specifying multiple parent folders (🅱️erhaps a csv, or inline csv list);
- Set up a releasing mechanism for people who know enough to download releases on Github but not how to build the code;
- ~~Look at setting up a real license~~ eh, maybe later

[^1]: Don't be a brainlet about this. What I am saying is: if you are a business that can afford to write a cheque to me with a sum of money, then you, the entity in a position to pay for this. Stop trying to use the legal system or accounting magic to shirk responsibility. And people wonder why others hate corporations.
