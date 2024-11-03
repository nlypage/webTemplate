import { S3 } from "https://deno.land/x/s3@0.5.0/mod.ts";
import { walk } from "@std/fs";
import { contentType } from "@std/media-types";
import { extname } from "@std/path/extname";

const s3 = new S3({
    accessKeyID: Deno.env.get("S3-ID")!,
    secretKey: Deno.env.get("S3-KEY")!,
    region: Deno.env.get("S3-REGION")!,
    endpointURL: Deno.env.get("S3-URL")!,
});

const bucket = s3.getBucket(Deno.env.get("S3-NAME")!);

for await (const entry of walk("../dist/")) {
    if (entry.isFile) {
        const filePath = entry.path;
        const fileRelativePath = filePath.replace(/\\/g, '/').replace("dist/", '');

        const body = await Deno.readFile(filePath);
        console.log(`uploading ${fileRelativePath}`)
        await bucket.putObject(fileRelativePath, body, {
            contentType: contentType(extname(filePath)) || 'application/octet-stream'
        }).catch(err => {
            console.log(`failed to upload ${filePath}. ${err}`);
        })
    }
}