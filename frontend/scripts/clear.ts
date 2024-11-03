import { S3 } from "https://deno.land/x/s3@0.5.0/mod.ts";

const s3 = new S3({
    accessKeyID: Deno.env.get("S3-ID")!,
    secretKey: Deno.env.get("S3-KEY")!,
    region: Deno.env.get("S3-REGION")!,
    endpointURL: Deno.env.get("S3-URL")!,
});

const bucket = s3.getBucket(Deno.env.get("S3-NAME")!);

const objects = bucket.listAllObjects({
    batchSize: 0
})
for await (const obj of objects) {
    console.log(`removing object ${obj.key}`)
    await bucket.deleteObject(obj.key!)
}