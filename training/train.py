# TensorFlow and tf.keras
import tensorflow as tf
from tensorflow import keras
import oss2

# Helper libraries
import os
import logging
import json
import tarfile
import traceback

akid = os.environ['ACCESS_KEY_ID']
akid_secret = os.environ['ACCESS_KEY_ID_SECRET']
region = os.environ['REGION']
bucket = os.environ['OSS_BUCKET']

logging.getLogger("tensorflow").setLevel(logging.INFO)
log = logging.getLogger('tensorflow')
log.info(tf.__version__)

fashion_mnist = keras.datasets.fashion_mnist
(train_images, train_labels), (test_images, test_labels) = fashion_mnist.load_data()

# scale the values to 0.0 to 1.0
train_images = train_images / 255.0
test_images = test_images / 255.0

# reshape for feeding into the model
train_images = train_images.reshape(train_images.shape[0], 28, 28, 1)
test_images = test_images.reshape(test_images.shape[0], 28, 28, 1)

class_names = ['T-shirt/top', 'Trouser', 'Pullover', 'Dress', 'Coat',
               'Sandal', 'Shirt', 'Sneaker', 'Bag', 'Ankle boot']

print('\ntrain_images.shape: {}, of {}'.format(train_images.shape, train_images.dtype))
print('test_images.shape: {}, of {}'.format(test_images.shape, test_images.dtype))

model = keras.Sequential([
    keras.layers.Conv2D(input_shape=(28,28,1), filters=8, kernel_size=3,
                        strides=2, activation='relu', name='Conv1'),
    keras.layers.Flatten(),
    keras.layers.Dense(10, activation=tf.nn.softmax, name='Softmax')
])
model.summary()

testing = False
epochs = 5

model.compile(optimizer=tf.optimizers.Adam(),
              loss='sparse_categorical_crossentropy',
              metrics=['accuracy'])
model.fit(train_images, train_labels, epochs=epochs)

test_loss, test_acc = model.evaluate(test_images, test_labels)
print('\nTest accuracy: {}'.format(test_acc))

# Fetch the Keras session and save the model
# The signature definition is defined by the input and output tensors,
# and stored with the default serving key

MODEL_DIR = "/tmp/tf/models"
version = 1
export_path = os.path.join(MODEL_DIR, str(version))
print('export_path = {}\n'.format(export_path))
if os.path.isdir(export_path):
    print('\nAlready saved a model, cleaning up\n')

tf.saved_model.save(
    model,
    export_path)

print('\nSaved model at %s' % export_path)

data = json.dumps({"signature_name": "serving_default", "instances": test_images[0:3].tolist()})
print('Data: {} ... {}'.format(data[:50], data[len(data)-52:]))

tarfile_name = "models.tar.gz"
tarfile_path = "/tmp/tf/" + tarfile_name
with tarfile.open(tarfile_path, "w:gz") as tar:
    tar.add(MODEL_DIR, arcname=os.path.basename(MODEL_DIR))

auth = oss2.Auth(akid, akid_secret)
oss_client = oss2.Bucket(auth, 'oss-%s.aliyuncs.com' % region, bucket)

try:
    model_file_oss_key = "fnf_k8s_demo_trained_models/models.tar.gz"
    print('\nUploading model %s to OSS %s %s' % (tarfile_path, bucket, model_file_oss_key))
    oss_client.put_object_from_file(model_file_oss_key, tarfile_path)
    print('\nUploaded model %s to OSS %s %s' % (tarfile_path, bucket, model_file_oss_key))
except:
    print(traceback.format_exc())

os.remove(tarfile_path)

