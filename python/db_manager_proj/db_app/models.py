from django.db import models

class Tiger(models.Model):
    """Biodata of a Tiger"""
    name = models.CharField(max_length = 200, db_column='name')
    dob = models.DateField(null=True,db_column='dob')

    class Meta:
        db_table = "tiger_bio"

class Spotting(models.Model):
    """Data regarding when the tiger is seen.
    For simplicity, we are using Decimal field to save cordinates, 
    for advanced developmenet we can use Point fields (PostGIS)"""
    tiger = models.ForeignKey(Tiger, on_delete=models.DO_NOTHING, related_name='tiger_spotting', default=None)
    seen_time = models.DateTimeField(db_column='seen_time')
    latitude = models.DecimalField(max_digits=22, decimal_places=16,db_column='latitude')
    longitude = models.DecimalField(max_digits=22, decimal_places=16, db_column='longitude')
    image = models.CharField(max_length = 250, db_column='image', default="")

    class Meta:
        db_table = "sighting_info"