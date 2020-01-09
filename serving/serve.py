import oss2

# Helper libraries
import os
import tarfile
import traceback
import subprocess

akid = os.environ['ACCESS_KEY_ID']
akid_secret = os.environ['ACCESS_KEY_ID_SECRET']
region = os.environ['REGION']
bucket = os.environ['OSS_BUCKET']

tarfile_name = "models.tar.gz"
tarfile_path = "/tmp/" + tarfile_name

auth = oss2.Auth(akid, akid_secret)
oss_client = oss2.Bucket(auth, 'oss-%s.aliyuncs.com' % region, bucket)
try:
    model_file_oss_key = "fnf_k8s_demo_trained_models/models.tar.gz"
    print('\nDownloading model from OSS %s %s to %s' % (bucket, model_file_oss_key, tarfile_path))
    oss_client.get_object_to_file(model_file_oss_key, tarfile_path)
    print('\nDownloaded model from OSS %s %s to %s' % (bucket, model_file_oss_key, tarfile_path))
except:
    print(traceback.format_exc())

tf = tarfile.open(tarfile_path)
tf.extractall()

subprocess.call(["tensorflow_model_server", "--rest_api_port=8501", "--model_name=fashion_model", "--model_base_path=/serving/models"])



